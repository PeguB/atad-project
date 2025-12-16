package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PeguB/atad-project/internal/models"
)

// CSVParser handles parsing CSV bank statements
type CSVParser struct {
	dateFormats []string
}

// NewCSVParser creates a new CSV parser with common date formats
func NewCSVParser() *CSVParser {
	return &CSVParser{
		dateFormats: []string{
			"2006-01-02",
			"02/01/2006",
			"01/02/2006",
			"2006/01/02",
			"Jan 02, 2006",
			"02-Jan-2006",
			"2006-01-02 15:04:05",
		},
	}
}

// ParseFile parses a CSV file and returns transactions
func (p *CSVParser) ParseFile(filename string) ([]*models.Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Detect column indices
	colMap := p.detectColumns(headers)
	if colMap["date"] == -1 || colMap["amount"] == -1 || colMap["description"] == -1 {
		return nil, fmt.Errorf("required columns not found (need: date, amount, description)")
	}

	var transactions []*models.Transaction
	lineNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNum, err)
		}
		lineNum++

		tx, err := p.parseRecord(record, colMap)
		if err != nil {
			// Skip invalid records with warning
			fmt.Printf("Warning: Skipping line %d: %v\n", lineNum, err)
			continue
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// detectColumns maps CSV headers to field indices
func (p *CSVParser) detectColumns(headers []string) map[string]int {
	colMap := map[string]int{
		"date":        -1,
		"description": -1,
		"amount":      -1,
		"category":    -1,
		"type":        -1,
	}

	for i, header := range headers {
		headerLower := strings.ToLower(strings.TrimSpace(header))

		// Date column
		if strings.Contains(headerLower, "date") || headerLower == "posted" || headerLower == "transaction date" {
			colMap["date"] = i
		}

		// Description column
		if strings.Contains(headerLower, "description") || strings.Contains(headerLower, "memo") ||
			strings.Contains(headerLower, "merchant") || headerLower == "payee" {
			colMap["description"] = i
		}

		// Amount column
		if strings.Contains(headerLower, "amount") || headerLower == "debit" ||
			headerLower == "credit" || headerLower == "transaction amount" {
			colMap["amount"] = i
		}

		// Category column (optional)
		if strings.Contains(headerLower, "category") || strings.Contains(headerLower, "type") {
			if colMap["category"] == -1 { // Only set if not already set
				colMap["category"] = i
			}
		}

		// Transaction type column (optional)
		if headerLower == "type" || headerLower == "transaction type" {
			colMap["type"] = i
		}
	}

	return colMap
}

// parseRecord converts a CSV record to a Transaction
func (p *CSVParser) parseRecord(record []string, colMap map[string]int) (*models.Transaction, error) {
	// Parse date
	dateStr := strings.TrimSpace(record[colMap["date"]])
	date, err := p.parseDate(dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date '%s': %w", dateStr, err)
	}

	// Parse amount
	amountStr := strings.TrimSpace(record[colMap["amount"]])
	amount, txType, err := p.parseAmount(amountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid amount '%s': %w", amountStr, err)
	}

	// Override type if explicitly specified
	if colMap["type"] != -1 && colMap["type"] < len(record) {
		typeStr := strings.ToLower(strings.TrimSpace(record[colMap["type"]]))
		if typeStr == "credit" || typeStr == "deposit" || typeStr == "income" {
			txType = "income"
		} else if typeStr == "debit" || typeStr == "withdrawal" || typeStr == "expense" {
			txType = "expense"
		}
	}

	// Get description
	description := strings.TrimSpace(record[colMap["description"]])
	if description == "" {
		description = "Imported transaction"
	}

	// Get category if available
	category := "Uncategorized"
	if colMap["category"] != -1 && colMap["category"] < len(record) {
		cat := strings.TrimSpace(record[colMap["category"]])
		if cat != "" {
			category = cat
		}
	}

	return &models.Transaction{
		Type:        txType,
		Amount:      amount,
		Category:    category,
		Description: description,
		Date:        date,
	}, nil
}

// parseDate tries multiple date formats
func (p *CSVParser) parseDate(dateStr string) (time.Time, error) {
	for _, format := range p.dateFormats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date format")
}

// parseAmount extracts amount and determines transaction type
func (p *CSVParser) parseAmount(amountStr string) (float64, string, error) {
	// Remove currency symbols and spaces
	cleaned := strings.TrimSpace(amountStr)
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "£", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	// Check for negative (expense) or positive (income)
	txType := "expense"
	if strings.HasPrefix(cleaned, "+") {
		txType = "income"
		cleaned = strings.TrimPrefix(cleaned, "+")
	} else if strings.HasPrefix(cleaned, "-") {
		cleaned = strings.TrimPrefix(cleaned, "-")
	} else if strings.HasSuffix(cleaned, "CR") || strings.HasSuffix(cleaned, "Cr") {
		// Credit notation
		txType = "income"
		cleaned = strings.TrimSuffix(cleaned, "CR")
		cleaned = strings.TrimSuffix(cleaned, "Cr")
	} else if strings.HasSuffix(cleaned, "DR") || strings.HasSuffix(cleaned, "Dr") {
		// Debit notation
		cleaned = strings.TrimSuffix(cleaned, "DR")
		cleaned = strings.TrimSuffix(cleaned, "Dr")
	}

	// Handle parentheses notation for negative amounts
	if strings.HasPrefix(cleaned, "(") && strings.HasSuffix(cleaned, ")") {
		cleaned = strings.TrimPrefix(cleaned, "(")
		cleaned = strings.TrimSuffix(cleaned, ")")
		txType = "expense"
	}

	cleaned = strings.TrimSpace(cleaned)

	amount, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, "", err
	}

	// Ensure amount is positive
	if amount < 0 {
		amount = -amount
	}

	return amount, txType, nil
}
