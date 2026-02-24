package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/repository"
	"github.com/homemenu/backend/service/generator"
)

type MealPlanService struct {
	planRepo     repository.MealPlanRepo
	itemRepo     repository.MealPlanItemRepo
	recipeRepo   repository.RecipeRepo
	ingredientRepo repository.IngredientRepo
	generator    generator.MenuGenerator
}

func NewMealPlanService(
	planRepo repository.MealPlanRepo,
	itemRepo repository.MealPlanItemRepo,
	recipeRepo repository.RecipeRepo,
	ingredientRepo repository.IngredientRepo,
	gen generator.MenuGenerator,
) *MealPlanService {
	return &MealPlanService{
		planRepo:     planRepo,
		itemRepo:     itemRepo,
		recipeRepo:   recipeRepo,
		ingredientRepo: ingredientRepo,
		generator:    gen,
	}
}

func (s *MealPlanService) Create(ctx context.Context, plan *model.MealPlan) error {
	plan.ShareToken = uuid.New().String()
	return s.planRepo.Create(ctx, plan)
}

func (s *MealPlanService) GetByID(ctx context.Context, id int64) (*model.MealPlan, error) {
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepo.ListByPlanID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.populateItemRecipes(ctx, items); err != nil {
		return nil, err
	}

	plan.Items = items
	return plan, nil
}

func (s *MealPlanService) GetByShareToken(ctx context.Context, token string) (*model.MealPlan, error) {
	plan, err := s.planRepo.GetByShareToken(ctx, token)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepo.ListByPlanID(ctx, plan.ID)
	if err != nil {
		return nil, err
	}

	if err := s.populateItemRecipes(ctx, items); err != nil {
		return nil, err
	}

	plan.Items = items
	return plan, nil
}

func (s *MealPlanService) List(ctx context.Context, userID int64) ([]model.MealPlan, error) {
	return s.planRepo.List(ctx, userID)
}

func (s *MealPlanService) Update(ctx context.Context, plan *model.MealPlan) error {
	return s.planRepo.Update(ctx, plan)
}

func (s *MealPlanService) Delete(ctx context.Context, id int64) error {
	if err := s.itemRepo.DeleteByPlanID(ctx, id); err != nil {
		return err
	}
	return s.planRepo.Delete(ctx, id)
}

func (s *MealPlanService) Generate(ctx context.Context, userID int64, config model.PlanConfig, startDate, endDate string) (*model.MealPlan, error) {
	// Get all recipes for this user
	recipes, err := s.recipeRepo.List(ctx, model.RecipeFilters{UserID: userID})
	if err != nil {
		return nil, err
	}

	// Load ingredients for filtering
	if len(recipes) > 0 {
		ids := make([]int64, len(recipes))
		for i, r := range recipes {
			ids[i] = r.ID
		}
		ingredientMap, err := s.ingredientRepo.ListByRecipeIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
		for i := range recipes {
			recipes[i].Ingredients = ingredientMap[recipes[i].ID]
		}
	}

	// Generate items using the engine
	items, err := s.generator.Generate(ctx, recipes, config, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Create the meal plan
	plan := &model.MealPlan{
		UserID:    userID,
		Name:      "菜单 " + startDate,
		StartDate: startDate,
		EndDate:   endDate,
		Config:    config,
	}
	if err := s.Create(ctx, plan); err != nil {
		return nil, err
	}

	// Set plan ID on items and create them
	for i := range items {
		items[i].MealPlanID = plan.ID
	}
	if err := s.itemRepo.BatchCreate(ctx, items); err != nil {
		return nil, err
	}

	// Reload with recipes populated
	return s.GetByID(ctx, plan.ID)
}

func (s *MealPlanService) UpdateItems(ctx context.Context, planID int64, items []model.MealPlanItem) error {
	if err := s.itemRepo.DeleteByPlanID(ctx, planID); err != nil {
		return err
	}
	for i := range items {
		items[i].MealPlanID = planID
	}
	return s.itemRepo.BatchCreate(ctx, items)
}

func (s *MealPlanService) RerollItem(ctx context.Context, planID, itemID int64, userID int64) (*model.MealPlanItem, error) {
	// Get current items to know what's already used
	currentItems, err := s.itemRepo.ListByPlanID(ctx, planID)
	if err != nil {
		return nil, err
	}

	usedRecipeIDs := make(map[int64]bool)
	var targetItem *model.MealPlanItem
	for i := range currentItems {
		usedRecipeIDs[currentItems[i].RecipeID] = true
		if currentItems[i].ID == itemID {
			targetItem = &currentItems[i]
		}
	}

	if targetItem == nil {
		return nil, ErrNotFound
	}

	// Get all recipes
	recipes, err := s.recipeRepo.List(ctx, model.RecipeFilters{UserID: userID})
	if err != nil {
		return nil, err
	}

	// Find a recipe not already in the plan
	for _, r := range recipes {
		if !usedRecipeIDs[r.ID] {
			targetItem.RecipeID = r.ID
			if err := s.itemRepo.Update(ctx, targetItem); err != nil {
				return nil, err
			}
			recipe, _ := s.recipeRepo.GetByID(ctx, r.ID)
			targetItem.Recipe = recipe
			return targetItem, nil
		}
	}

	return targetItem, nil
}

func (s *MealPlanService) populateItemRecipes(ctx context.Context, items []model.MealPlanItem) error {
	if len(items) == 0 {
		return nil
	}

	recipeIDs := make([]int64, 0, len(items))
	seen := make(map[int64]bool)
	for _, item := range items {
		if !seen[item.RecipeID] {
			recipeIDs = append(recipeIDs, item.RecipeID)
			seen[item.RecipeID] = true
		}
	}

	recipes, err := s.recipeRepo.ListByIDs(ctx, recipeIDs)
	if err != nil {
		return err
	}

	// Load ingredients
	ingredientMap, err := s.ingredientRepo.ListByRecipeIDs(ctx, recipeIDs)
	if err != nil {
		return err
	}

	recipeMap := make(map[int64]*model.Recipe)
	for i := range recipes {
		recipes[i].Ingredients = ingredientMap[recipes[i].ID]
		if recipes[i].Ingredients == nil {
			recipes[i].Ingredients = []model.Ingredient{}
		}
		recipeMap[recipes[i].ID] = &recipes[i]
	}

	for i := range items {
		if r, ok := recipeMap[items[i].RecipeID]; ok {
			items[i].Recipe = r
		}
	}

	return nil
}
