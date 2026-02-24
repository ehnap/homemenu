package model

import "time"

type MealPlan struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"user_id"`
	Name       string         `json:"name"`
	StartDate  string         `json:"start_date"`
	EndDate    string         `json:"end_date"`
	Config     PlanConfig     `json:"config"`
	ShareToken string         `json:"share_token,omitempty"`
	Items      []MealPlanItem `json:"items,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type PlanConfig struct {
	MealTypes       []string `json:"meal_types"`
	DishesPerMeal   map[string]DishCount `json:"dishes_per_meal,omitempty"`
	TastePreference string   `json:"taste_preference,omitempty"`
	PreferIngredients []string `json:"prefer_ingredients,omitempty"`
	ExcludeIngredients []string `json:"exclude_ingredients,omitempty"`
	UseAI           bool     `json:"use_ai,omitempty"`
}

type DishCount struct {
	Meat      int `json:"meat"`
	Vegetable int `json:"vegetable"`
	Soup      int `json:"soup"`
}
