package service

import (
	"context"

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
		if err := s.ingredientRepo.BatchCreate(ctx, recipe.ID, recipe.Ingredients); err != nil {
			return err
		}
	}

	return nil
}

func (s *RecipeService) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	recipe, err := s.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ingredients, err := s.ingredientRepo.ListByRecipeID(ctx, id)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients = ingredients

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
			recipes[i].Ingredients = ingredientMap[recipes[i].ID]
			if recipes[i].Ingredients == nil {
				recipes[i].Ingredients = []model.Ingredient{}
			}
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
		if err := s.ingredientRepo.BatchCreate(ctx, recipe.ID, recipe.Ingredients); err != nil {
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
