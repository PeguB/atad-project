package repository

import (
	"database/sql"
	"fmt"

	"github.com/PeguB/atad-project/internal/models"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create adds a new budget
func (r *BudgetRepository) Create(budget *models.Budget) error {
	query := `
		INSERT INTO budgets (category, amount, period, start_date, end_date)
		VALUES (?, ?, ?, ?, ?)
	`

	var startDate, endDate interface{}
	if !budget.StartDate.IsZero() {
		startDate = budget.StartDate
	}
	if !budget.EndDate.IsZero() {
		endDate = budget.EndDate
	}

	result, err := r.db.Exec(query, budget.Category, budget.Amount, budget.Period, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to create budget: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	budget.ID = id
	return nil
}

// GetByCategory retrieves a budget for a specific category
func (r *BudgetRepository) GetByCategory(category string) (*models.Budget, error) {
	query := `
		SELECT id, category, amount, period
		FROM budgets
		WHERE category = ?
	`

	budget := &models.Budget{}
	err := r.db.QueryRow(query, category).Scan(
		&budget.ID,
		&budget.Category,
		&budget.Amount,
		&budget.Period,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No budget set for this category
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	return budget, nil
}

// GetByCategoryAndPeriod retrieves a budget for a specific category and period
func (r *BudgetRepository) GetByCategoryAndPeriod(category, period string) (*models.Budget, error) {
	query := `
		SELECT id, category, amount, period, start_date, end_date
		FROM budgets
		WHERE category = ? AND period = ?
	`

	budget := &models.Budget{}
	var startDate, endDate sql.NullTime
	err := r.db.QueryRow(query, category, period).Scan(
		&budget.ID,
		&budget.Category,
		&budget.Amount,
		&budget.Period,
		&startDate,
		&endDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No budget set for this category and period
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	if startDate.Valid {
		budget.StartDate = startDate.Time
	}
	if endDate.Valid {
		budget.EndDate = endDate.Time
	}

	return budget, nil
}

// GetByCategoryAndDateRange retrieves a budget for a specific category and custom date range
func (r *BudgetRepository) GetByCategoryAndDateRange(category string, startDate, endDate interface{}) (*models.Budget, error) {
	query := `
		SELECT id, category, amount, period, start_date, end_date
		FROM budgets
		WHERE category = ? AND period = 'custom' AND start_date = ? AND end_date = ?
	`

	budget := &models.Budget{}
	var sd, ed sql.NullTime
	err := r.db.QueryRow(query, category, startDate, endDate).Scan(
		&budget.ID,
		&budget.Category,
		&budget.Amount,
		&budget.Period,
		&sd,
		&ed,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	if sd.Valid {
		budget.StartDate = sd.Time
	}
	if ed.Valid {
		budget.EndDate = ed.Time
	}

	return budget, nil
}

// GetAll retrieves all budgets
func (r *BudgetRepository) GetAll() ([]*models.Budget, error) {
	query := `
		SELECT id, category, amount, period, start_date, end_date
		FROM budgets
		ORDER BY category
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query budgets: %w", err)
	}
	defer rows.Close()

	var budgets []*models.Budget
	for rows.Next() {
		budget := &models.Budget{}
		var startDate, endDate sql.NullTime
		err := rows.Scan(&budget.ID, &budget.Category, &budget.Amount, &budget.Period, &startDate, &endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		if startDate.Valid {
			budget.StartDate = startDate.Time
		}
		if endDate.Valid {
			budget.EndDate = endDate.Time
		}
		budgets = append(budgets, budget)
	}

	return budgets, rows.Err()
}

// Update modifies an existing budget
func (r *BudgetRepository) Update(budget *models.Budget) error {
	query := `
		UPDATE budgets
		SET amount = ?
		WHERE category = ? AND period = ?
	`

	result, err := r.db.Exec(query, budget.Amount, budget.Category, budget.Period)
	if err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("budget not found")
	}

	return nil
}

// Delete removes a budget by category and period
func (r *BudgetRepository) Delete(category, period string) error {
	query := `DELETE FROM budgets WHERE category = ? AND period = ?`
	result, err := r.db.Exec(query, category, period)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("budget not found")
	}

	return nil
}

// GetSpending calculates total spending for a category in a given date range
func (r *BudgetRepository) GetSpending(category string, startDate, endDate interface{}) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE category = ? 
		AND type = 'expense'
		AND date >= ? AND date <= ?
	`

	var total float64
	err := r.db.QueryRow(query, category, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get spending: %w", err)
	}

	return total, nil
}
