package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/PeguB/atad-project/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create adds a new transaction
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (date, description, amount, category, type, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		tx.Date,
		tx.Description,
		tx.Amount,
		tx.Category,
		tx.Type,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tx.ID = id
	return nil
}

// GetAll retrieves all transactions
func (r *TransactionRepository) GetAll() ([]*models.Transaction, error) {
	query := `
		SELECT id, date, description, amount, category, type, created_at
		FROM transactions
		ORDER BY date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(&tx.ID, &tx.Date, &tx.Description, &tx.Amount, &tx.Category, &tx.Type, &tx.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// GetByType retrieves transactions by type (income or expense)
func (r *TransactionRepository) GetByType(txType string) ([]*models.Transaction, error) {
	query := `
		SELECT id, date, description, amount, category, type, created_at
		FROM transactions
		WHERE type = ?
		ORDER BY date DESC
	`

	rows, err := r.db.Query(query, txType)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by type: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(&tx.ID, &tx.Date, &tx.Description, &tx.Amount, &tx.Category, &tx.Type, &tx.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// Delete removes a transaction
func (r *TransactionRepository) Delete(id int64) error {
	query := `DELETE FROM transactions WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
