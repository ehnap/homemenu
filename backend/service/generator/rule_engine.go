package generator

import (
	"context"
	"math/rand"
	"time"

	"github.com/homemenu/backend/model"
)

type RuleEngine struct{}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (e *RuleEngine) Generate(ctx context.Context, recipes []model.Recipe, config model.PlanConfig, startDate, endDate string) ([]model.MealPlanItem, error) {
	if len(recipes) == 0 {
		return []model.MealPlanItem{}, nil
	}

	// Filter recipes based on config
	filtered := filterRecipes(recipes, config)
	if len(filtered) == 0 {
		filtered = recipes
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	mealTypes := config.MealTypes
	if len(mealTypes) == 0 {
		mealTypes = []string{"lunch", "dinner"}
	}

	var items []model.MealPlanItem
	usedRecipes := make(map[int64]bool)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		for _, mealType := range mealTypes {
			dishCount := getDishCount(config, mealType)
			for i := 0; i < dishCount; i++ {
				recipe := pickRecipe(filtered, usedRecipes, rng)
				if recipe == nil {
					// Reset used recipes if we've exhausted all
					usedRecipes = make(map[int64]bool)
					recipe = pickRecipe(filtered, usedRecipes, rng)
				}
				if recipe == nil {
					continue
				}
				usedRecipes[recipe.ID] = true
				items = append(items, model.MealPlanItem{
					RecipeID:  recipe.ID,
					Date:      dateStr,
					MealType:  mealType,
					SortOrder: i,
				})
			}
		}
	}

	return items, nil
}

func filterRecipes(recipes []model.Recipe, config model.PlanConfig) []model.Recipe {
	var result []model.Recipe
	for _, r := range recipes {
		// Check taste preference
		if config.TastePreference != "" {
			hasTag := false
			for _, tag := range r.Tags {
				if tag == config.TastePreference {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		// Check excluded ingredients
		if len(config.ExcludeIngredients) > 0 {
			excluded := false
			for _, ing := range r.Ingredients {
				for _, ex := range config.ExcludeIngredients {
					if ing.Name == ex {
						excluded = true
						break
					}
				}
				if excluded {
					break
				}
			}
			if excluded {
				continue
			}
		}

		result = append(result, r)
	}
	return result
}

func pickRecipe(recipes []model.Recipe, used map[int64]bool, rng *rand.Rand) *model.Recipe {
	available := make([]int, 0)
	for i, r := range recipes {
		if !used[r.ID] {
			available = append(available, i)
		}
	}
	if len(available) == 0 {
		return nil
	}
	idx := available[rng.Intn(len(available))]
	return &recipes[idx]
}

func getDishCount(config model.PlanConfig, mealType string) int {
	if config.DishesPerMeal != nil {
		if dc, ok := config.DishesPerMeal[mealType]; ok {
			total := dc.Meat + dc.Vegetable + dc.Soup
			if total > 0 {
				return total
			}
		}
	}
	// Defaults
	switch mealType {
	case "breakfast":
		return 1
	case "lunch":
		return 3
	case "dinner":
		return 2
	default:
		return 2
	}
}
