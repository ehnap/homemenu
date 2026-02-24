package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/service"
)

type MealPlanHandler struct {
	mealPlanService *service.MealPlanService
	shoppingService *service.ShoppingService
}

func NewMealPlanHandler(mealPlanService *service.MealPlanService, shoppingService *service.ShoppingService) *MealPlanHandler {
	return &MealPlanHandler{
		mealPlanService: mealPlanService,
		shoppingService: shoppingService,
	}
}

func (h *MealPlanHandler) List(c *gin.Context) {
	plans, err := h.mealPlanService.List(c.Request.Context(), GetUserID(c))
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, plans)
}

func (h *MealPlanHandler) Create(c *gin.Context) {
	var plan model.MealPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	plan.UserID = GetUserID(c)

	if err := h.mealPlanService.Create(c.Request.Context(), &plan); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, plan)
}

func (h *MealPlanHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	plan, err := h.mealPlanService.GetByID(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusNotFound, "meal plan not found")
		return
	}
	Success(c, plan)
}

func (h *MealPlanHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var plan model.MealPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	plan.ID = id
	plan.UserID = GetUserID(c)

	if err := h.mealPlanService.Update(c.Request.Context(), &plan); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, plan)
}

func (h *MealPlanHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.mealPlanService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, nil)
}

func (h *MealPlanHandler) Generate(c *gin.Context) {
	var req struct {
		StartDate string           `json:"start_date" binding:"required"`
		EndDate   string           `json:"end_date" binding:"required"`
		Config    model.PlanConfig `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	plan, err := h.mealPlanService.Generate(c.Request.Context(), GetUserID(c), req.Config, req.StartDate, req.EndDate)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, plan)
}

func (h *MealPlanHandler) UpdateItems(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Items []model.MealPlanItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.mealPlanService.UpdateItems(c.Request.Context(), id, req.Items); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	plan, _ := h.mealPlanService.GetByID(c.Request.Context(), id)
	Success(c, plan)
}

func (h *MealPlanHandler) RerollItem(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid plan id")
		return
	}
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.mealPlanService.RerollItem(c.Request.Context(), planID, itemID, GetUserID(c))
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, item)
}

func (h *MealPlanHandler) ShoppingList(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	mode := c.DefaultQuery("mode", "weekly")
	list, err := h.shoppingService.GetShoppingList(c.Request.Context(), id, mode)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, list)
}
