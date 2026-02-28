package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/homemenu/backend/model"
)

type IngredientRepo struct {
	db *sql.DB
}

func NewIngredientRepo(db *sql.DB) *IngredientRepo {
	return &IngredientRepo{db: db}
}

func (r *IngredientRepo) BatchCreate(ctx context.Context, recipeID int64, ingredients []model.Ingredient) error {
	if len(ingredients) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO recipe_ingredients (recipe_id, name, amount, unit) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, ing := range ingredients {
		result, err := stmt.ExecContext(ctx, recipeID, ing.Name, ing.Amount, ing.Unit)
		if err != nil {
			return err
		}
		id, _ := result.LastInsertId()
		ing.ID = id
	}

	return tx.Commit()
}

func (r *IngredientRepo) ListByRecipeID(ctx context.Context, recipeID int64) ([]model.Ingredient, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, recipe_id, name, amount, unit FROM recipe_ingredients WHERE recipe_id = ?", recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []model.Ingredient
	for rows.Next() {
		var ing model.Ingredient
		if err := rows.Scan(&ing.ID, &ing.RecipeID, &ing.Name, &ing.Amount, &ing.Unit); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}
	if ingredients == nil {
		ingredients = []model.Ingredient{}
	}
	return ingredients, nil
}

func (r *IngredientRepo) ListByRecipeIDs(ctx context.Context, recipeIDs []int64) (map[int64][]model.Ingredient, error) {
	result := make(map[int64][]model.Ingredient)
	if len(recipeIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(recipeIDs))
	args := make([]interface{}, len(recipeIDs))
	for i, id := range recipeIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(
		"SELECT id, recipe_id, name, amount, unit FROM recipe_ingredients WHERE recipe_id IN (%s)",
		strings.Join(placeholders, ","),
	)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ing model.Ingredient
		if err := rows.Scan(&ing.ID, &ing.RecipeID, &ing.Name, &ing.Amount, &ing.Unit); err != nil {
			return nil, err
		}
		result[ing.RecipeID] = append(result[ing.RecipeID], ing)
	}
	return result, nil
}

func (r *IngredientRepo) DeleteByRecipeID(ctx context.Context, recipeID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM recipe_ingredients WHERE recipe_id = ?", recipeID)
	return err
}

func (r *IngredientRepo) SuggestNames(ctx context.Context, query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.db.QueryContext(ctx,
		"SELECT DISTINCT name FROM recipe_ingredients WHERE name LIKE ? ORDER BY name LIMIT ?",
		"%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	if names == nil {
		names = []string{}
	}
	return names, nil
}
