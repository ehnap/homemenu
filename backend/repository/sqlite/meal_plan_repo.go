package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/homemenu/backend/model"
)

type MealPlanRepo struct {
	db *sql.DB
}

func NewMealPlanRepo(db *sql.DB) *MealPlanRepo {
	return &MealPlanRepo{db: db}
}

func (r *MealPlanRepo) Create(ctx context.Context, plan *model.MealPlan) error {
	configJSON, _ := json.Marshal(plan.Config)
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO meal_plans (user_id, name, start_date, end_date, config, share_token)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		plan.UserID, plan.Name, plan.StartDate, plan.EndDate, string(configJSON), plan.ShareToken,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	plan.ID = id
	return nil
}

func (r *MealPlanRepo) GetByID(ctx context.Context, id int64) (*model.MealPlan, error) {
	plan := &model.MealPlan{}
	var configJSON string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, start_date, end_date, config, share_token, created_at, updated_at
		 FROM meal_plans WHERE id = ?`, id,
	).Scan(&plan.ID, &plan.UserID, &plan.Name, &plan.StartDate, &plan.EndDate,
		&configJSON, &plan.ShareToken, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(configJSON), &plan.Config)
	return plan, nil
}

func (r *MealPlanRepo) GetByShareToken(ctx context.Context, token string) (*model.MealPlan, error) {
	plan := &model.MealPlan{}
	var configJSON string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, start_date, end_date, config, share_token, created_at, updated_at
		 FROM meal_plans WHERE share_token = ?`, token,
	).Scan(&plan.ID, &plan.UserID, &plan.Name, &plan.StartDate, &plan.EndDate,
		&configJSON, &plan.ShareToken, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(configJSON), &plan.Config)
	return plan, nil
}

func (r *MealPlanRepo) List(ctx context.Context, userID int64) ([]model.MealPlan, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, name, start_date, end_date, config, share_token, created_at, updated_at
		 FROM meal_plans WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []model.MealPlan
	for rows.Next() {
		var plan model.MealPlan
		var configJSON string
		if err := rows.Scan(&plan.ID, &plan.UserID, &plan.Name, &plan.StartDate, &plan.EndDate,
			&configJSON, &plan.ShareToken, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(configJSON), &plan.Config)
		plans = append(plans, plan)
	}
	if plans == nil {
		plans = []model.MealPlan{}
	}
	return plans, nil
}

func (r *MealPlanRepo) Update(ctx context.Context, plan *model.MealPlan) error {
	configJSON, _ := json.Marshal(plan.Config)
	_, err := r.db.ExecContext(ctx,
		`UPDATE meal_plans SET name=?, start_date=?, end_date=?, config=?, updated_at=CURRENT_TIMESTAMP
		 WHERE id=?`,
		plan.Name, plan.StartDate, plan.EndDate, string(configJSON), plan.ID,
	)
	return err
}

func (r *MealPlanRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM meal_plans WHERE id = ?", id)
	return err
}

// MealPlanItemRepo

type MealPlanItemRepo struct {
	db *sql.DB
}

func NewMealPlanItemRepo(db *sql.DB) *MealPlanItemRepo {
	return &MealPlanItemRepo{db: db}
}

func (r *MealPlanItemRepo) BatchCreate(ctx context.Context, items []model.MealPlanItem) error {
	if len(items) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO meal_plan_items (meal_plan_id, recipe_id, date, meal_type, sort_order) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range items {
		result, err := stmt.ExecContext(ctx, items[i].MealPlanID, items[i].RecipeID, items[i].Date, items[i].MealType, items[i].SortOrder)
		if err != nil {
			return err
		}
		id, _ := result.LastInsertId()
		items[i].ID = id
	}
	return tx.Commit()
}

func (r *MealPlanItemRepo) ListByPlanID(ctx context.Context, planID int64) ([]model.MealPlanItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, meal_plan_id, recipe_id, date, meal_type, sort_order
		 FROM meal_plan_items WHERE meal_plan_id = ? ORDER BY date, meal_type, sort_order`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.MealPlanItem
	for rows.Next() {
		var item model.MealPlanItem
		if err := rows.Scan(&item.ID, &item.MealPlanID, &item.RecipeID, &item.Date, &item.MealType, &item.SortOrder); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []model.MealPlanItem{}
	}
	return items, nil
}

func (r *MealPlanItemRepo) DeleteByPlanID(ctx context.Context, planID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM meal_plan_items WHERE meal_plan_id = ?", planID)
	return err
}

func (r *MealPlanItemRepo) DeleteByID(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM meal_plan_items WHERE id = ?", id)
	return err
}

func (r *MealPlanItemRepo) Update(ctx context.Context, item *model.MealPlanItem) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE meal_plan_items SET recipe_id=?, date=?, meal_type=?, sort_order=? WHERE id=?",
		item.RecipeID, item.Date, item.MealType, item.SortOrder, item.ID,
	)
	return err
}
