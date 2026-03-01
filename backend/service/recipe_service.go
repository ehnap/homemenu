package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/homemenu/backend/model"
	"github.com/homemenu/backend/repository"
)

type RecipeService struct {
	recipeRepo     repository.RecipeRepo
	ingredientRepo repository.IngredientRepo
}

func NewRecipeService(recipeRepo repository.RecipeRepo, ingredientRepo repository.IngredientRepo) *RecipeService {
	return &RecipeService{recipeRepo: recipeRepo, ingredientRepo: ingredientRepo}
}

func (s *RecipeService) Create(ctx context.Context, recipe *model.Recipe) error {
	if err := s.recipeRepo.Create(ctx, recipe); err != nil {
		return err
	}

	if len(recipe.Ingredients) > 0 {
		if err := s.ingredientRepo.BatchCreate(ctx, recipe.ID, recipe.Ingredients, "ingredient"); err != nil {
			return err
		}
	}

	if len(recipe.Seasonings) > 0 {
		if err := s.ingredientRepo.BatchCreate(ctx, recipe.ID, recipe.Seasonings, "seasoning"); err != nil {
			return err
		}
	}

	return nil
}

// splitIngredients separates a flat ingredient list into ingredients and seasonings by Category.
func splitIngredients(all []model.Ingredient) (ingredients, seasonings []model.Ingredient) {
	for _, ing := range all {
		if ing.Category == "seasoning" {
			seasonings = append(seasonings, ing)
		} else {
			ingredients = append(ingredients, ing)
		}
	}
	if ingredients == nil {
		ingredients = []model.Ingredient{}
	}
	if seasonings == nil {
		seasonings = []model.Ingredient{}
	}
	return
}

func (s *RecipeService) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	recipe, err := s.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	all, err := s.ingredientRepo.ListByRecipeID(ctx, id)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients, recipe.Seasonings = splitIngredients(all)

	return recipe, nil
}

func (s *RecipeService) List(ctx context.Context, filters model.RecipeFilters) ([]model.Recipe, error) {
	recipes, err := s.recipeRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

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
			all := ingredientMap[recipes[i].ID]
			if all == nil {
				all = []model.Ingredient{}
			}
			recipes[i].Ingredients, recipes[i].Seasonings = splitIngredients(all)
		}
	}

	return recipes, nil
}

func (s *RecipeService) Update(ctx context.Context, recipe *model.Recipe) error {
	if err := s.recipeRepo.Update(ctx, recipe); err != nil {
		return err
	}

	if err := s.ingredientRepo.DeleteByRecipeID(ctx, recipe.ID); err != nil {
		return err
	}

	if len(recipe.Ingredients) > 0 {
		if err := s.ingredientRepo.BatchCreate(ctx, recipe.ID, recipe.Ingredients, "ingredient"); err != nil {
			return err
		}
	}

	if len(recipe.Seasonings) > 0 {
		if err := s.ingredientRepo.BatchCreate(ctx, recipe.ID, recipe.Seasonings, "seasoning"); err != nil {
			return err
		}
	}

	return nil
}

func (s *RecipeService) Delete(ctx context.Context, id int64) error {
	if err := s.ingredientRepo.DeleteByRecipeID(ctx, id); err != nil {
		return err
	}
	return s.recipeRepo.Delete(ctx, id)
}

func (s *RecipeService) SuggestIngredients(ctx context.Context, query string, limit int) ([]string, error) {
	return s.ingredientRepo.SuggestNames(ctx, query, limit)
}

func (s *RecipeService) GenerateShareToken(ctx context.Context, id int64) (string, error) {
	recipe, err := s.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("recipe not found: %w", err)
	}

	if recipe.ShareToken != "" {
		return recipe.ShareToken, nil
	}

	token := uuid.New().String()
	recipe.ShareToken = token
	if err := s.recipeRepo.Update(ctx, recipe); err != nil {
		return "", err
	}
	return token, nil
}

func (s *RecipeService) GetByShareToken(ctx context.Context, token string) (*model.Recipe, error) {
	recipe, err := s.recipeRepo.GetByShareToken(ctx, token)
	if err != nil {
		return nil, err
	}

	all, err := s.ingredientRepo.ListByRecipeID(ctx, recipe.ID)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients, recipe.Seasonings = splitIngredients(all)

	return recipe, nil
}
