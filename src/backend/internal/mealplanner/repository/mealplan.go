package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/platepilot/backend/internal/mealplanner/domain"
)

// ErrMealPlanNotFound is returned when no plan exists for the week.
var ErrMealPlanNotFound = errors.New("meal plan not found")

// GetWeekPlan returns the saved week plan for a user and start date.
func (r *Repository) GetWeekPlan(ctx context.Context, userID uuid.UUID, startDate time.Time) (*domain.WeekPlan, error) {
	var planID uuid.UUID
	var dbStartDate time.Time
	var endDate time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT id, start_date, end_date
		FROM meal_plans
		WHERE user_id = $1 AND start_date = $2
	`, userID, startDate).Scan(&planID, &dbStartDate, &endDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMealPlanNotFound
		}
		return nil, fmt.Errorf("get meal plan: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT s.slot_date, s.meal_type, s.recipe_id, r.name, r.description
		FROM meal_plan_slots s
		LEFT JOIN recipes r ON r.id = s.recipe_id
		WHERE s.plan_id = $1
		ORDER BY s.slot_date, s.meal_type
	`, planID)
	if err != nil {
		return nil, fmt.Errorf("list meal plan slots: %w", err)
	}
	defer rows.Close()

	slots := make([]domain.MealSlot, 0)
	for rows.Next() {
		var slotDate time.Time
		var mealType string
		var recipeID uuid.UUID
		var recipeName *string
		var recipeDescription *string
		if err := rows.Scan(&slotDate, &mealType, &recipeID, &recipeName, &recipeDescription); err != nil {
			return nil, fmt.Errorf("scan meal plan slot: %w", err)
		}

		slot := domain.MealSlot{
			Date:     slotDate,
			MealType: mealType,
			RecipeID: recipeID,
		}
		if recipeName != nil {
			slot.RecipeName = *recipeName
		}
		if recipeDescription != nil {
			slot.RecipeDescription = *recipeDescription
		}
		slots = append(slots, slot)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate meal plan slots: %w", rows.Err())
	}

	return &domain.WeekPlan{
		UserID:    userID,
		StartDate: dbStartDate,
		EndDate:   endDate,
		Slots:     slots,
	}, nil
}

// UpsertWeekPlan creates or updates a week plan and its slots.
func (r *Repository) UpsertWeekPlan(ctx context.Context, plan domain.WeekPlan) (*domain.WeekPlan, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var planID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO meal_plans (user_id, start_date, end_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, start_date)
		DO UPDATE SET end_date = EXCLUDED.end_date, updated_at = NOW()
		RETURNING id
	`, plan.UserID, plan.StartDate, plan.EndDate).Scan(&planID)
	if err != nil {
		return nil, fmt.Errorf("upsert meal plan: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM meal_plan_slots WHERE plan_id = $1`, planID)
	if err != nil {
		return nil, fmt.Errorf("clear meal plan slots: %w", err)
	}

	for _, slot := range plan.Slots {
		if slot.RecipeID == uuid.Nil {
			continue
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO meal_plan_slots (plan_id, slot_date, meal_type, recipe_id)
			VALUES ($1, $2, $3, $4)
		`, planID, slot.Date, slot.MealType, slot.RecipeID)
		if err != nil {
			return nil, fmt.Errorf("insert meal plan slot: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit meal plan: %w", err)
	}

	return r.GetWeekPlan(ctx, plan.UserID, plan.StartDate)
}
