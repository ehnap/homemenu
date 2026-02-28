package service

import (
	"context"

	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/repository"
)

type SettingsService struct {
	settingsRepo repository.SettingsRepo
}

func NewSettingsService(settingsRepo repository.SettingsRepo) *SettingsService {
	return &SettingsService{settingsRepo: settingsRepo}
}

var llmKeys = []string{"llm_base_url", "llm_api_key", "llm_model"}

func (s *SettingsService) GetLLMSettings(ctx context.Context, userID int64) (*model.LLMSettings, error) {
	values, err := s.settingsRepo.GetMulti(ctx, userID, llmKeys)
	if err != nil {
		return nil, err
	}
	return &model.LLMSettings{
		BaseURL: values["llm_base_url"],
		APIKey:  values["llm_api_key"],
		Model:   values["llm_model"],
	}, nil
}

func (s *SettingsService) UpdateLLMSettings(ctx context.Context, userID int64, settings *model.LLMSettings) error {
	values := map[string]string{
		"llm_base_url": settings.BaseURL,
		"llm_api_key":  settings.APIKey,
		"llm_model":    settings.Model,
	}
	return s.settingsRepo.SetMulti(ctx, userID, values)
}
