package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/service"
)

type RecipeHandler struct {
	recipeService *service.RecipeService
}

func NewRecipeHandler(recipeService *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

func (h *RecipeHandler) List(c *gin.Context) {
	filters := model.RecipeFilters{
		UserID:     GetUserID(c),
		Tag:        c.Query("tag"),
		Ingredient: c.Query("ingredient"),
		Query:      c.Query("q"),
		Difficulty: c.Query("difficulty"),
	}

	// Support multiple ingredients: ?ingredients=菠菜,猪肉
	if ings := c.Query("ingredients"); ings != "" {
		filters.Ingredients = strings.Split(ings, ",")
	}

	recipes, err := h.recipeService.List(c.Request.Context(), filters)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, recipes)
}

func (h *RecipeHandler) Create(c *gin.Context) {
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	recipe.UserID = GetUserID(c)

	if recipe.Name == "" {
		Error(c, http.StatusBadRequest, "recipe name is required")
		return
	}

	if err := h.recipeService.Create(c.Request.Context(), &recipe); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, recipe)
}

func (h *RecipeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	recipe, err := h.recipeService.GetByID(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusNotFound, "recipe not found")
		return
	}

	Success(c, recipe)
}

func (h *RecipeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	recipe.ID = id
	recipe.UserID = GetUserID(c)

	if err := h.recipeService.Update(c.Request.Context(), &recipe); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, recipe)
}

func (h *RecipeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.recipeService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, nil)
}
