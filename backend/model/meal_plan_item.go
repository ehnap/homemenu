package model

type MealPlanItem struct {
	ID         int64   `json:"id"`
	MealPlanID int64   `json:"meal_plan_id"`
	RecipeID   int64   `json:"recipe_id"`
	Date       string  `json:"date"`
	MealType   string  `json:"meal_type"`
	SortOrder  int     `json:"sort_order"`
	Recipe     *Recipe `json:"recipe,omitempty"`
}
