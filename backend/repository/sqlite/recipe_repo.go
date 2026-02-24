package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/homemenu/backend/model"
)

type RecipeRepo struct {
	db *sql.DB
}

func NewRecipeRepo(db *sql.DB) *RecipeRepo {
	return &RecipeRepo{db: db}
}

func (r *RecipeRepo) Create(ctx context.Context, recipe *model.Recipe) error {
	stepsJSON, _ := json.Marshal(recipe.Steps)
	tagsJSON, _ := json.Marshal(recipe.Tags)

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO recipes (user_id, name, steps, cook_time, difficulty, tags, cover_image, calories, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		recipe.UserID, recipe.Name, string(stepsJSON), recipe.CookTime,
		recipe.Difficulty, string(tagsJSON), recipe.CoverImage, recipe.Calories, recipe.Notes,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	recipe.ID = id
	return nil
}

func (r *RecipeRepo) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	recipe := &model.Recipe{}
	var stepsJSON, tagsJSON string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, steps, cook_time, difficulty, tags, cover_image, calories, notes, created_at, updated_at
		 FROM recipes WHERE id = ?`, id,
	).Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &stepsJSON, &recipe.CookTime,
		&recipe.Difficulty, &tagsJSON, &recipe.CoverImage, &recipe.Calories, &recipe.Notes,
		&recipe.CreatedAt, &recipe.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(stepsJSON), &recipe.Steps)
	json.Unmarshal([]byte(tagsJSON), &recipe.Tags)
	if recipe.Steps == nil {
		recipe.Steps = []model.Step{}
	}
	if recipe.Tags == nil {
		recipe.Tags = []string{}
	}
	return recipe, nil
}

func (r *RecipeRepo) List(ctx context.Context, filters model.RecipeFilters) ([]model.Recipe, error) {
	query := "SELECT DISTINCT r.id, r.user_id, r.name, r.steps, r.cook_time, r.difficulty, r.tags, r.cover_image, r.calories, r.notes, r.created_at, r.updated_at FROM recipes r"
	var conditions []string
	var args []interface{}

	if filters.UserID > 0 {
		conditions = append(conditions, "r.user_id = ?")
		args = append(args, filters.UserID)
	}

	if filters.Tag != "" {
		conditions = append(conditions, "r.tags LIKE ?")
		args = append(args, "%\""+filters.Tag+"\"%")
	}

	if filters.Difficulty != "" {
		conditions = append(conditions, "r.difficulty = ?")
		args = append(args, filters.Difficulty)
	}

	if filters.Query != "" {
		conditions = append(conditions, "r.name LIKE ?")
		args = append(args, "%"+filters.Query+"%")
	}

	// Single ingredient filter
	if filters.Ingredient != "" {
		query += " JOIN recipe_ingredients ri ON r.id = ri.recipe_id"
		conditions = append(conditions, "ri.name LIKE ?")
		args = append(args, "%"+filters.Ingredient+"%")
	}

	// Multiple ingredients filter (INTERSECT approach)
	if len(filters.Ingredients) > 0 && filters.Ingredient == "" {
		var subQueries []string
		for _, ing := range filters.Ingredients {
			subQueries = append(subQueries, "SELECT recipe_id FROM recipe_ingredients WHERE name LIKE ?")
			args = append(args, "%"+ing+"%")
		}
		// We need r.id to be in the intersection of all ingredient queries
		intersect := strings.Join(subQueries, " INTERSECT ")
		conditions = append(conditions, fmt.Sprintf("r.id IN (%s)", intersect))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY r.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []model.Recipe
	for rows.Next() {
		var recipe model.Recipe
		var stepsJSON, tagsJSON string
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &stepsJSON, &recipe.CookTime,
			&recipe.Difficulty, &tagsJSON, &recipe.CoverImage, &recipe.Calories, &recipe.Notes,
			&recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(stepsJSON), &recipe.Steps)
		json.Unmarshal([]byte(tagsJSON), &recipe.Tags)
		if recipe.Steps == nil {
			recipe.Steps = []model.Step{}
		}
		if recipe.Tags == nil {
			recipe.Tags = []string{}
		}
		recipes = append(recipes, recipe)
	}
	if recipes == nil {
		recipes = []model.Recipe{}
	}
	return recipes, nil
}

func (r *RecipeRepo) Update(ctx context.Context, recipe *model.Recipe) error {
	stepsJSON, _ := json.Marshal(recipe.Steps)
	tagsJSON, _ := json.Marshal(recipe.Tags)

	_, err := r.db.ExecContext(ctx,
		`UPDATE recipes SET name=?, steps=?, cook_time=?, difficulty=?, tags=?, cover_image=?, calories=?, notes=?, updated_at=CURRENT_TIMESTAMP
		 WHERE id=?`,
		recipe.Name, string(stepsJSON), recipe.CookTime, recipe.Difficulty,
		string(tagsJSON), recipe.CoverImage, recipe.Calories, recipe.Notes, recipe.ID,
	)
	return err
}

func (r *RecipeRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM recipes WHERE id = ?", id)
	return err
}

func (r *RecipeRepo) ListByIDs(ctx context.Context, ids []int64) ([]model.Recipe, error) {
	if len(ids) == 0 {
		return []model.Recipe{}, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(
		"SELECT id, user_id, name, steps, cook_time, difficulty, tags, cover_image, calories, notes, created_at, updated_at FROM recipes WHERE id IN (%s)",
		strings.Join(placeholders, ","),
	)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []model.Recipe
	for rows.Next() {
		var recipe model.Recipe
		var stepsJSON, tagsJSON string
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &stepsJSON, &recipe.CookTime,
			&recipe.Difficulty, &tagsJSON, &recipe.CoverImage, &recipe.Calories, &recipe.Notes,
			&recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(stepsJSON), &recipe.Steps)
		json.Unmarshal([]byte(tagsJSON), &recipe.Tags)
		if recipe.Steps == nil {
			recipe.Steps = []model.Step{}
		}
		if recipe.Tags == nil {
			recipe.Tags = []string{}
		}
		recipes = append(recipes, recipe)
	}
	if recipes == nil {
		recipes = []model.Recipe{}
	}
	return recipes, nil
}
