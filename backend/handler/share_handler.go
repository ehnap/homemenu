package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homemenu/backend/service"
)

type ShareHandler struct {
	mealPlanService *service.MealPlanService
	shoppingService *service.ShoppingService
	recipeService   *service.RecipeService
}

func NewShareHandler(mealPlanService *service.MealPlanService, shoppingService *service.ShoppingService, recipeService *service.RecipeService) *ShareHandler {
	return &ShareHandler{
		mealPlanService: mealPlanService,
		shoppingService: shoppingService,
		recipeService:   recipeService,
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

func (h *ShareHandler) GetRecipeByToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		Error(c, http.StatusBadRequest, "missing share token")
		return
	}

	recipe, err := h.recipeService.GetByShareToken(c.Request.Context(), token)
	if err != nil {
		Error(c, http.StatusNotFound, "shared recipe not found")
		return
	}

	Success(c, recipe)
}
