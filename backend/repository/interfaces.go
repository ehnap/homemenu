package repository

import (
	"context"

	"github.com/homemenu/backend/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type RecipeRepo interface {
	Create(ctx context.Context, recipe *model.Recipe) error
	GetByID(ctx context.Context, id int64) (*model.Recipe, error)
	List(ctx context.Context, filters model.RecipeFilters) ([]model.Recipe, error)
	Update(ctx context.Context, recipe *model.Recipe) error
	Delete(ctx context.Context, id int64) error
	ListByIDs(ctx context.Context, ids []int64) ([]model.Recipe, error)
}

type IngredientRepo interface {
	BatchCreate(ctx context.Context, recipeID int64, ingredients []model.Ingredient) error
	ListByRecipeID(ctx context.Context, recipeID int64) ([]model.Ingredient, error)
	ListByRecipeIDs(ctx context.Context, recipeIDs []int64) (map[int64][]model.Ingredient, error)
	DeleteByRecipeID(ctx context.Context, recipeID int64) error
}

type MealPlanRepo interface {
	Create(ctx context.Context, plan *model.MealPlan) error
	GetByID(ctx context.Context, id int64) (*model.MealPlan, error)
	GetByShareToken(ctx context.Context, token string) (*model.MealPlan, error)
	List(ctx context.Context, userID int64) ([]model.MealPlan, error)
	Update(ctx context.Context, plan *model.MealPlan) error
	Delete(ctx context.Context, id int64) error
}

type MealPlanItemRepo interface {
	BatchCreate(ctx context.Context, items []model.MealPlanItem) error
	ListByPlanID(ctx context.Context, planID int64) ([]model.MealPlanItem, error)
	DeleteByPlanID(ctx context.Context, planID int64) error
	DeleteByID(ctx context.Context, id int64) error
	Update(ctx context.Context, item *model.MealPlanItem) error
}
