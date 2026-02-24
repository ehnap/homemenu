package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homemenu/backend/service"
)

type ShareHandler struct {
	mealPlanService *service.MealPlanService
	shoppingService *service.ShoppingService
}

func NewShareHandler(mealPlanService *service.MealPlanService, shoppingService *service.ShoppingService) *ShareHandler {
	return &ShareHandler{
		mealPlanService: mealPlanService,
		shoppingService: shoppingService,
	}
}

func (h *ShareHandler) GetByToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		Error(c, http.StatusBadRequest, "missing share token")
		return
	}

	plan, err := h.mealPlanService.GetByShareToken(c.Request.Context(), token)
	if err != nil {
		Error(c, http.StatusNotFound, "shared meal plan not found")
		return
	}

	// Also include shopping list
	shoppingList, _ := h.shoppingService.GetShoppingList(c.Request.Context(), plan.ID, "daily")

	Success(c, gin.H{
		"meal_plan":     plan,
		"shopping_list": shoppingList,
	})
}
