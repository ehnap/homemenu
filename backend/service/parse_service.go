package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/homemenu/backend/config"
	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/service/llm"
)

type ParseService struct {
	llmClient       *llm.Client
	settingsService *SettingsService
	defaultLLMCfg   config.LLMConfig
}

func NewParseService(llmClient *llm.Client, settingsService *SettingsService, defaultCfg config.LLMConfig) *ParseService {
	return &ParseService{
		llmClient:       llmClient,
		settingsService: settingsService,
		defaultLLMCfg:   defaultCfg,
	}
}

func (s *ParseService) getLLMConfig(ctx context.Context, userID int64) (llm.Config, error) {
	// Try user settings first
	userSettings, err := s.settingsService.GetLLMSettings(ctx, userID)
	if err == nil && userSettings.APIKey != "" && userSettings.BaseURL != "" {
		return llm.Config{
			BaseURL: userSettings.BaseURL,
			APIKey:  userSettings.APIKey,
			Model:   userSettings.Model,
		}, nil
	}

	// Fall back to config.yaml
	if s.defaultLLMCfg.APIKey != "" && s.defaultLLMCfg.BaseURL != "" {
		return llm.Config{
			BaseURL: s.defaultLLMCfg.BaseURL,
			APIKey:  s.defaultLLMCfg.APIKey,
			Model:   s.defaultLLMCfg.Model,
		}, nil
	}

	return llm.Config{}, fmt.Errorf("LLM 未配置，请在设置页面配置 API Key")
}

func (s *ParseService) ParseRecipeText(ctx context.Context, userID int64, text string) (*model.Recipe, error) {
	cfg, err := s.getLLMConfig(ctx, userID)
	if err != nil {
		return nil, err
	}

	prompt := `你是一个菜谱解析助手。请从以下文本中提取菜谱信息，返回 JSON 格式。

要求：
1. 提取菜名（name）
2. 提取食材列表（ingredients），每个食材包含 name、amount、unit
3. 提取做法步骤（steps），每个步骤包含 order 和 description
4. 尽量识别标签（tags），如口味（咸、辣、清淡、甜、酸）、荤素、菜系等
5. 如果能识别难度（difficulty：简单/中等/复杂），也请提取
6. 如果能识别烹饪时间（cook_time，单位分钟），也请提取

返回格式：
{
  "name": "菜名",
  "ingredients": [{"name": "食材名", "amount": "用量", "unit": "单位"}],
  "steps": [{"order": 1, "description": "步骤描述"}],
  "tags": ["标签1", "标签2"],
  "difficulty": "简单",
  "cook_time": 30
}

只返回 JSON，不要其他文字。

以下是需要解析的菜谱文本：
` + text

	messages := []llm.Message{
		{Role: "user", Content: prompt},
	}

	content, err := s.llmClient.Chat(ctx, cfg, messages, 0.3)
	if err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	// Clean markdown code block wrappers
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var recipe model.Recipe
	if err := json.Unmarshal([]byte(content), &recipe); err != nil {
		return nil, fmt.Errorf("解析结果格式错误: %w", err)
	}

	return &recipe, nil
}
