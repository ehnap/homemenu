package model

type Ingredient struct {
	ID       int64  `json:"id"`
	RecipeID int64  `json:"recipe_id"`
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Unit     string `json:"unit"`
}
