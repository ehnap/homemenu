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
	parseService  *service.ParseService
}

func NewRecipeHandler(recipeService *service.RecipeService, parseService *service.ParseService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService, parseService: parseService}
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

func (h *RecipeHandler) ParseText(c *gin.Context) {
	var req struct {
		Text string `json:"text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if len(req.Text) == 0 {
		Error(c, http.StatusBadRequest, "text is required")
		return
	}

	if len(req.Text) > 10000 {
		Error(c, http.StatusBadRequest, "text too long (max 10000 characters)")
		return
	}

	userID := GetUserID(c)
	recipe, err := h.parseService.ParseRecipeText(c.Request.Context(), userID, req.Text)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, recipe)
}

func (h *RecipeHandler) SuggestIngredients(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		Success(c, []string{})
		return
	}

	names, err := h.recipeService.SuggestIngredients(c.Request.Context(), q, 10)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, names)
}

func (h *RecipeHandler) GenerateShareToken(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	token, err := h.recipeService.GenerateShareToken(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, gin.H{"share_token": token})
}
