package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/service"
)

type SettingsHandler struct {
	settingsService *service.SettingsService
}

func NewSettingsHandler(settingsService *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService}
}

func (h *SettingsHandler) GetLLMSettings(c *gin.Context) {
	userID := GetUserID(c)

	settings, err := h.settingsService.GetLLMSettings(c.Request.Context(), userID)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Mask API key: show only last 4 characters
	if len(settings.APIKey) > 4 {
		settings.APIKey = strings.Repeat("*", len(settings.APIKey)-4) + settings.APIKey[len(settings.APIKey)-4:]
	}

	Success(c, settings)
}

func (h *SettingsHandler) UpdateLLMSettings(c *gin.Context) {
	userID := GetUserID(c)

	var req model.LLMSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	// If the API key looks masked (contains consecutive *), don't overwrite the stored value
	if strings.Contains(req.APIKey, "***") {
		existing, err := h.settingsService.GetLLMSettings(c.Request.Context(), userID)
		if err == nil && existing.APIKey != "" {
			req.APIKey = existing.APIKey
		}
	}

	if err := h.settingsService.UpdateLLMSettings(c.Request.Context(), userID, &req); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, nil)
}
