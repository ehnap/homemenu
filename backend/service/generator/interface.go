package generator

import (
	"context"

	"github.com/homemenu/backend/model"
)

type MenuGenerator interface {
	Generate(ctx context.Context, recipes []model.Recipe, config model.PlanConfig, startDate, endDate string) ([]model.MealPlanItem, error)
}
