package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/repository"
)

type ShoppingService struct {
	itemRepo       repository.MealPlanItemRepo
	recipeRepo     repository.RecipeRepo
	ingredientRepo repository.IngredientRepo
}

func NewShoppingService(
	itemRepo repository.MealPlanItemRepo,
	recipeRepo repository.RecipeRepo,
	ingredientRepo repository.IngredientRepo,
) *ShoppingService {
	return &ShoppingService{
		itemRepo:       itemRepo,
		recipeRepo:     recipeRepo,
		ingredientRepo: ingredientRepo,
	}
}

type ShoppingItem struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

type DailyShoppingList struct {
	Date  string         `json:"date"`
	Items []ShoppingItem `json:"items"`
}

type ShoppingList struct {
	Daily  []DailyShoppingList `json:"daily,omitempty"`
	Weekly []ShoppingItem      `json:"weekly,omitempty"`
}

func (s *ShoppingService) GetShoppingList(ctx context.Context, planID int64, mode string) (*ShoppingList, error) {
	items, err := s.itemRepo.ListByPlanID(ctx, planID)
	if err != nil {
		return nil, err
	}

	// Collect all recipe IDs
	recipeIDs := make([]int64, 0)
	seen := make(map[int64]bool)
	for _, item := range items {
		if !seen[item.RecipeID] {
			recipeIDs = append(recipeIDs, item.RecipeID)
			seen[item.RecipeID] = true
		}
	}

	// Get ingredients for all recipes
	ingredientMap, err := s.ingredientRepo.ListByRecipeIDs(ctx, recipeIDs)
	if err != nil {
		return nil, err
	}

	result := &ShoppingList{}

	if mode == "daily" {
		// Group by date
		dateRecipes := make(map[string][]int64)
		dateOrder := make([]string, 0)
		for _, item := range items {
			if _, exists := dateRecipes[item.Date]; !exists {
				dateOrder = append(dateOrder, item.Date)
			}
			dateRecipes[item.Date] = append(dateRecipes[item.Date], item.RecipeID)
		}

		for _, date := range dateOrder {
			daily := DailyShoppingList{Date: date}
			merged := mergeIngredients(dateRecipes[date], ingredientMap)
			daily.Items = merged
			result.Daily = append(result.Daily, daily)
		}
	} else {
		// Weekly: merge all
		allRecipeIDs := make([]int64, 0)
		for _, item := range items {
			allRecipeIDs = append(allRecipeIDs, item.RecipeID)
		}
		result.Weekly = mergeIngredients(allRecipeIDs, ingredientMap)
	}

	return result, nil
}

func mergeIngredients(recipeIDs []int64, ingredientMap map[int64][]model.Ingredient) []ShoppingItem {
	type key struct {
		name string
		unit string
	}
	merged := make(map[key]float64)
	order := make([]key, 0)

	for _, rid := range recipeIDs {
		for _, ing := range ingredientMap[rid] {
			k := key{name: ing.Name, unit: ing.Unit}
			if _, exists := merged[k]; !exists {
				order = append(order, k)
			}
			amount, err := strconv.ParseFloat(ing.Amount, 64)
			if err == nil {
				merged[k] += amount
			} else {
				// Non-numeric amount, just keep it
				merged[k] = 0
			}
		}
	}

	result := make([]ShoppingItem, 0, len(order))
	for _, k := range order {
		item := ShoppingItem{Name: k.name, Unit: k.unit}
		if merged[k] > 0 {
			item.Amount = fmt.Sprintf("%.0f", merged[k])
		}
		result = append(result, item)
	}
	return result
}
