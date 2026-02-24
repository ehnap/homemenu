package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/homemenu/backend/model"
)

type AIEngine struct {
	baseURL  string
	apiKey   string
	model    string
	fallback *RuleEngine
}

func NewAIEngine(baseURL, apiKey, modelName string) *AIEngine {
	return &AIEngine{
		baseURL:  baseURL,
		apiKey:   apiKey,
		model:    modelName,
		fallback: NewRuleEngine(),
	}
}

func (e *AIEngine) Generate(ctx context.Context, recipes []model.Recipe, config model.PlanConfig, startDate, endDate string) ([]model.MealPlanItem, error) {
	if e.apiKey == "" || e.baseURL == "" {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	// Build recipe summary for LLM
	recipeSummary := make([]map[string]interface{}, 0, len(recipes))
	for _, r := range recipes {
		recipeSummary = append(recipeSummary, map[string]interface{}{
			"id":         r.ID,
			"name":       r.Name,
			"tags":       r.Tags,
			"difficulty": r.Difficulty,
			"cook_time":  r.CookTime,
		})
	}

	prompt := fmt.Sprintf(`你是一个家庭菜单规划助手。请根据以下菜谱库和约束条件，生成从 %s 到 %s 的菜单。

菜谱库：
%s

约束条件：
- 餐次：%v
- 口味偏好：%s
- 优先食材：%v
- 排除食材：%v

请返回 JSON 数组格式，每个元素包含 recipe_id, date, meal_type, sort_order。
只返回 JSON，不要其他文字。`,
		startDate, endDate,
		mustJSON(recipeSummary),
		config.MealTypes,
		config.TastePreference,
		config.PreferIngredients,
		config.ExcludeIngredients,
	)

	// Call OpenAI-compatible API
	reqBody := map[string]interface{}{
		"model": e.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", e.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	// Parse OpenAI response
	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil || len(apiResp.Choices) == 0 {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	// Parse the LLM's JSON response into items
	var items []model.MealPlanItem
	content := apiResp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &items); err != nil {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	// Validate recipe IDs
	validIDs := make(map[int64]bool)
	for _, r := range recipes {
		validIDs[r.ID] = true
	}
	var validItems []model.MealPlanItem
	for _, item := range items {
		if validIDs[item.RecipeID] {
			validItems = append(validItems, item)
		}
	}

	if len(validItems) == 0 {
		return e.fallback.Generate(ctx, recipes, config, startDate, endDate)
	}

	return validItems, nil
}

func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
