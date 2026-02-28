package model

import "time"

type Recipe struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	Name        string       `json:"name"`
	Steps       []Step       `json:"steps"`
	CookTime    int          `json:"cook_time"`
	Difficulty  string       `json:"difficulty"`
	Tags        []string     `json:"tags"`
	CoverImage  string       `json:"cover_image"`
	Calories    int          `json:"calories"`
	Notes       string       `json:"notes"`
	Ingredients []Ingredient `json:"ingredients"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Step struct {
	Order       int    `json:"order"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url,omitempty"`
}

type RecipeFilters struct {
	UserID      int64
	Tag         string
	Ingredient  string
	Ingredients []string
	Query       string
	Difficulty  string
}
