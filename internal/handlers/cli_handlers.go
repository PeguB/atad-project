package handlers

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PeguB/atad-project/internal/database"
	"github.com/PeguB/atad-project/internal/models"
	"github.com/PeguB/atad-project/internal/parser"
	"github.com/PeguB/atad-project/internal/repository"
	"github.com/PeguB/atad-project/internal/service"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CommandHandler defines the interface for all command handlers
type CommandHandler interface {
	Handle()
}

// CLIHandler manages all CLI command operations
type CLIHandler struct {
	db              *database.Database
	txRepo          *repository.TransactionRepository
	budgetRepo      *repository.BudgetRepository
	categoryService *service.CategoryService
}

// NewCLIHandler creates a new CLI handler instance
func NewCLIHandler() *CLIHandler {
	return &CLIHandler{}
}

// InitDatabase initializes the database connection and repositories
func (h *CLIHandler) InitDatabase() error {
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	h.db = db
	h.txRepo = repository.NewTransactionRepository(db.DB)
	h.budgetRepo = repository.NewBudgetRepository(db.DB)
	h.categoryService = service.NewCategoryService()
	return nil
}

// Close closes the database connection
func (h *CLIHandler) Close() {
	if h.db != nil {
		h.db.Close()
	}
}

// AddCommand handles the 'add' subcommand
type AddCommand struct {
	Handler *CLIHandler
}

func (c *AddCommand) Handle() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	txType := addCmd.String("type", "", "Transaction type: income or expense (required)")
	description := addCmd.String("desc", "", "Transaction description (required)")
	amount := addCmd.Float64("amount", 0, "Transaction amount (required)")
	category := addCmd.String("category", "", "Transaction category (optional, auto-categorized if not provided)")
	date := addCmd.String("date", "", "Transaction date in DD/MM/YYYY format (optional, defaults to today)")

	addCmd.Parse(os.Args[2:])

	// Validate required fields
	if *txType == "" || *description == "" || *amount == 0 {
		fmt.Println("Error: -type, -desc, and -amount are required")
		fmt.Println("\nUsage: atad add -type <income|expense> -desc <description> -amount <amount> [-category <category>] [-date <DD/MM/YYYY>]")
		fmt.Println("\nExample: atad add -type expense -desc \"Grocery shopping\" -amount 75.50 -category Groceries")
		os.Exit(1)
	}

	if *txType != "income" && *txType != "expense" {
		fmt.Println("Error: -type must be either 'income' or 'expense'")
		os.Exit(1)
	}

	// Initialize database
	if err := c.Handler.InitDatabase(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer c.Handler.Close()

	// Auto-categorize if no category provided
	finalCategory := *category
	if finalCategory == "" {
		finalCategory = c.Handler.categoryService.CategorizeTransaction(*description)
		fmt.Printf("Auto-categorized as: %s\n", finalCategory)
	}

	// Parse date
	var txDate time.Time
	var err error
	if *date == "" {
		txDate = time.Now()
	} else {
		txDate, err = time.Parse("02/01/2006", *date)
		if err != nil {
			fmt.Printf("Error: Invalid date format. Use DD/MM/YYYY (e.g., 15/12/2025)\n")
			os.Exit(1)
		}
	}

	// Create transaction
	tx := &models.Transaction{
		Date:        txDate,
		Description: *description,
		Amount:      *amount,
		Category:    finalCategory,
		Type:        *txType,
	}

	err = c.Handler.txRepo.Create(tx)
	if err != nil {
		fmt.Printf("Error saving transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Transaction added successfully!\n")
	caser := cases.Title(language.English)
	fmt.Printf("   Type: %s\n", caser.String(*txType))
	fmt.Printf("   Description: %s\n", *description)
	fmt.Printf("   Amount: $%.2f\n", *amount)
	fmt.Printf("   Category: %s\n", finalCategory)
	fmt.Printf("   Date: %s\n", txDate.Format("02/01/2006"))

	// Check budget
	if *txType == "expense" {
		budget, _ := c.Handler.budgetRepo.GetByCategory(finalCategory)
		if budget != nil {
			spending, _ := c.Handler.budgetRepo.GetSpending(finalCategory, budget.StartDate, budget.EndDate)
			percentUsed := (spending / budget.Amount) * 100
			if spending > budget.Amount {
				fmt.Printf("\n⚠️  Over budget! Spent: $%.2f / $%.2f (%.0f%%)\n", spending, budget.Amount, percentUsed)
			} else if percentUsed >= 80 {
				fmt.Printf("\n⚠️  Budget warning: $%.2f / $%.2f (%.0f%%)\n", spending, budget.Amount, percentUsed)
			}
		}
	}
}

// ListCommand handles the 'list' subcommand
type ListCommand struct {
	Handler *CLIHandler
}

func (c *ListCommand) Handle() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	txType := listCmd.String("type", "all", "Filter by type: all, income, or expense")
	limit := listCmd.Int("limit", 20, "Number of transactions to display")

	listCmd.Parse(os.Args[2:])

	if err := c.Handler.InitDatabase(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer c.Handler.Close()

	var transactions []*models.Transaction
	var err error
	if *txType == "all" {
		transactions, err = c.Handler.txRepo.GetAll()
	} else {
		transactions, err = c.Handler.txRepo.GetByType(*txType)
	}

	if err != nil {
		fmt.Printf("Error retrieving transactions: %v\n", err)
		os.Exit(1)
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions found.")
		return
	}

	// Limit results
	displayCount := len(transactions)
	if *limit > 0 && *limit < displayCount {
		displayCount = *limit
	}

	fmt.Printf("\n📋 Transactions (%d of %d)\n", displayCount, len(transactions))
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("%-12s %-10s %-25s %-15s %10s\n", "Date", "Type", "Description", "Category", "Amount")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")

	caser := cases.Title(language.English)
	for i := 0; i < displayCount; i++ {
		tx := transactions[i]
		typeIcon := "💰"
		if tx.Type == "expense" {
			typeIcon = "💸"
		}
		fmt.Printf("%-12s %-10s %-25s %-15s %10.2f\n",
			tx.Date.Format("02/01/2006"),
			typeIcon+" "+caser.String(tx.Type),
			TruncateString(tx.Description, 25),
			TruncateString(tx.Category, 15),
			tx.Amount)
	}

	if len(transactions) > displayCount {
		fmt.Printf("\n... and %d more. Use -limit flag to see more.\n", len(transactions)-displayCount)
	}
}

// ReportCommand handles the 'report' subcommand
type ReportCommand struct {
	Handler *CLIHandler
}

func (c *ReportCommand) Handle() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: atad report <income|expense> [-period <all|month|year>]")
		os.Exit(1)
	}

	reportType := os.Args[2]
	reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
	period := reportCmd.String("period", "month", "Time period: all, month, or year")

	reportCmd.Parse(os.Args[3:])

	if err := c.Handler.InitDatabase(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer c.Handler.Close()

	var startDate, endDate time.Time
	now := time.Now()

	switch *period {
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), 12, 31, 0, 0, 0, 0, now.Location())
	case "all":
		// No date filter
	default:
		fmt.Println("Error: -period must be 'all', 'month', or 'year'")
		os.Exit(1)
	}

	transactions, err := c.Handler.txRepo.GetAll()
	if err != nil {
		fmt.Printf("Error retrieving transactions: %v\n", err)
		os.Exit(1)
	}

	total := 0.0
	byCategory := make(map[string]float64)

	for _, tx := range transactions {
		if tx.Type != reportType {
			continue
		}

		if *period != "all" {
			if tx.Date.Before(startDate) || tx.Date.After(endDate) {
				continue
			}
		}

		total += tx.Amount
		byCategory[tx.Category] += tx.Amount
	}

	caser := cases.Title(language.English)
	periodName := caser.String(*period)
	if *period == "month" {
		periodName = now.Format("January 2006")
	} else if *period == "year" {
		periodName = strconv.Itoa(now.Year())
	}

	fmt.Printf("\n📊 %s Report - %s\n\n", caser.String(reportType), periodName)
	fmt.Printf("Total %s: $%.2f\n\n", caser.String(reportType), total)

	if len(byCategory) > 0 {
		// Draw bar chart
		DrawCategoryBarChart(byCategory, total)

		fmt.Println("\nBreakdown by Category:")
		fmt.Println("─────────────────────────────────────")

		// Sort categories by amount for consistent display
		type categoryAmount struct {
			category string
			amount   float64
		}
		var sorted []categoryAmount
		for cat, amt := range byCategory {
			sorted = append(sorted, categoryAmount{cat, amt})
		}
		// Sort descending by amount
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j].amount > sorted[i].amount {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}

		for _, item := range sorted {
			percentage := (item.amount / total) * 100
			// Get color for this category
			colorStyle := GetCategoryColor(item.category)
			colorIndicator := colorStyle.Render("█")
			fmt.Printf("  %s %-20s $%8.2f  (%.1f%%)\n", colorIndicator, item.category, item.amount, percentage)
		}
	} else {
		fmt.Printf("No %s transactions found for this period.\n", reportType)
	}
}

// BudgetCommand handles the 'budget' subcommand
type BudgetCommand struct {
	Handler *CLIHandler
}

func (c *BudgetCommand) Handle() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  atad budget list                                      # List all budgets")
		fmt.Println("  atad budget set <category> <amount> <start> <end>     # Set a budget")
		fmt.Println("  atad budget check <category>                          # Check budget status")
		os.Exit(1)
	}

	action := os.Args[2]

	if err := c.Handler.InitDatabase(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer c.Handler.Close()

	switch action {
	case "list":
		c.handleList()
	case "set":
		c.handleSet()
	case "check":
		c.handleCheck()
	default:
		fmt.Printf("Unknown budget action: %s\n", action)
		os.Exit(1)
	}
}

func (c *BudgetCommand) handleList() {
	budgets, err := c.Handler.budgetRepo.GetAll()
	if err != nil {
		fmt.Printf("Error retrieving budgets: %v\n", err)
		os.Exit(1)
	}

	if len(budgets) == 0 {
		fmt.Println("No budgets set.")
		return
	}

	fmt.Println("\n💰 Budgets")
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Printf("%-20s %12s  %-24s\n", "Category", "Amount", "Period")
	fmt.Println("─────────────────────────────────────────────────────────────")

	for _, budget := range budgets {
		period := fmt.Sprintf("%s - %s",
			budget.StartDate.Format("02/01/2006"),
			budget.EndDate.Format("02/01/2006"))
		fmt.Printf("%-20s $%11.2f  %-24s\n", budget.Category, budget.Amount, period)
	}
}

func (c *BudgetCommand) handleSet() {
	if len(os.Args) < 7 {
		fmt.Println("Usage: atad budget set <category> <amount> <start_date> <end_date>")
		fmt.Println("Example: atad budget set Groceries 500 01/12/2025 31/12/2025")
		os.Exit(1)
	}

	category := os.Args[3]
	amount, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Println("Error: Invalid amount")
		os.Exit(1)
	}

	startDate, err := time.Parse("02/01/2006", os.Args[5])
	if err != nil {
		fmt.Println("Error: Invalid start date. Use DD/MM/YYYY format")
		os.Exit(1)
	}

	endDate, err := time.Parse("02/01/2006", os.Args[6])
	if err != nil {
		fmt.Println("Error: Invalid end date. Use DD/MM/YYYY format")
		os.Exit(1)
	}

	budget := &models.Budget{
		Category:  category,
		Amount:    amount,
		Period:    "custom",
		StartDate: startDate,
		EndDate:   endDate,
	}

	err = c.Handler.budgetRepo.Create(budget)
	if err != nil {
		fmt.Printf("Error creating budget: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Budget set successfully!\n")
	fmt.Printf("   Category: %s\n", category)
	fmt.Printf("   Amount: $%.2f\n", amount)
	fmt.Printf("   Period: %s to %s\n", startDate.Format("02/01/2006"), endDate.Format("02/01/2006"))
}

func (c *BudgetCommand) handleCheck() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: atad budget check <category>")
		os.Exit(1)
	}

	category := os.Args[3]
	budget, err := c.Handler.budgetRepo.GetByCategory(category)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if budget == nil {
		fmt.Printf("No budget set for category '%s'\n", category)
		return
	}

	spending, err := c.Handler.budgetRepo.GetSpending(category, budget.StartDate, budget.EndDate)
	if err != nil {
		fmt.Printf("Error calculating spending: %v\n", err)
		os.Exit(1)
	}

	percentUsed := (spending / budget.Amount) * 100
	remaining := budget.Amount - spending

	fmt.Printf("\n💰 Budget Status: %s\n", category)
	fmt.Println("─────────────────────────────────────")
	fmt.Printf("Budget:     $%.2f\n", budget.Amount)
	fmt.Printf("Spent:      $%.2f (%.0f%%)\n", spending, percentUsed)
	fmt.Printf("Remaining:  $%.2f\n", remaining)
	fmt.Printf("Period:     %s - %s\n",
		budget.StartDate.Format("02/01/2006"),
		budget.EndDate.Format("02/01/2006"))

	if spending > budget.Amount {
		fmt.Printf("\n⚠️  Over budget by $%.2f!\n", spending-budget.Amount)
	} else if percentUsed >= 80 {
		fmt.Println("\n⚠️  Warning: 80% or more of budget used")
	} else {
		fmt.Println("\n✅ Within budget")
	}
}

// SearchCommand handles the 'search' subcommand
type SearchCommand struct {
	Handler *CLIHandler
}

func (c *SearchCommand) Handle() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: atad search <query>")
		fmt.Println("Example: atad search \"coffee\"")
		os.Exit(1)
	}

	query := strings.ToLower(os.Args[2])

	if err := c.Handler.InitDatabase(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer c.Handler.Close()

	transactions, err := c.Handler.txRepo.GetAll()
	if err != nil {
		fmt.Printf("Error retrieving transactions: %v\n", err)
		os.Exit(1)
	}

	var results []*models.Transaction
	for _, tx := range transactions {
		if strings.Contains(strings.ToLower(tx.Description), query) ||
			strings.Contains(strings.ToLower(tx.Category), query) {
			results = append(results, tx)
		}
	}

	if len(results) == 0 {
		fmt.Printf("No transactions found matching '%s'\n", os.Args[2])
		return
	}

	fmt.Printf("\n🔍 Search Results for '%s' (%d found)\n", os.Args[2], len(results))
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("%-12s %-10s %-25s %-15s %10s\n", "Date", "Type", "Description", "Category", "Amount")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")

	caser := cases.Title(language.English)
	for _, tx := range results {
		typeIcon := "💰"
		if tx.Type == "expense" {
			typeIcon = "💸"
		}
		fmt.Printf("%-12s %-10s %-25s %-15s %10.2f\n",
			tx.Date.Format("02/01/2006"),
			typeIcon+" "+caser.String(tx.Type),
			TruncateString(tx.Description, 25),
			TruncateString(tx.Category, 15),
			tx.Amount)
	}
}
// ImportCommand handles the 'import' subcommand
type ImportCommand struct {
	Handler *CLIHandler
}

func (c *ImportCommand) Handle() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: atad import <file.csv> [--auto-categorize] [--skip-duplicates]")
		fmt.Println("\nSupported format: CSV (comma-separated values)")
		fmt.Println("\nOptions:")
		fmt.Println("  --auto-categorize    Automatically categorize imported transactions")
		fmt.Println("  --skip-duplicates    Skip transactions that appear to be duplicates")
		fmt.Println("\nExample:")
		fmt.Println("  atad import statement.csv --auto-categorize")
		os.Exit(1)
	}

	filename := os.Args[2]
	autoCategorize := false
	skipDuplicates := false

	// Parse flags
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--auto-categorize":
			autoCategorize = true
		case "--skip-duplicates":
			skipDuplicates = true
		}
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found\n", filename)
		os.Exit(1)
	}

	fmt.Printf("📥 Importing transactions from %s...\n\n", filename)

	// Parse CSV file
	csvParser := parser.NewCSVParser()
	transactions, err := csvParser.ParseFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions found in file")
		return
	}

	fmt.Printf("Found %d transactions\n", len(transactions))

	// Auto-categorize if requested
	if autoCategorize {
		fmt.Println("Applying automatic categorization...")
		for _, tx := range transactions {
			if tx.Category == "Uncategorized" || tx.Category == "" {
				suggestedCat := c.Handler.categoryService.CategorizeTransaction(tx.Description)
				if suggestedCat != "" {
					tx.Category = suggestedCat
				}
			}
		}
	}

	// Remove duplicates if requested
	imported := 0
	skipped := 0
	errors := 0

	if skipDuplicates {
		fmt.Println("Checking for duplicates...")
	}

	for _, tx := range transactions {
		// Check for duplicates
		if skipDuplicates {
			isDuplicate, err := c.Handler.txRepo.IsDuplicate(tx)
			if err != nil {
				fmt.Printf("Warning: Error checking duplicate: %v\n", err)
			} else if isDuplicate {
				skipped++
				continue
			}
		}

		// Save transaction
		err := c.Handler.txRepo.Create(tx)
		if err != nil {
			fmt.Printf("Error importing transaction: %v\n", err)
			errors++
		} else {
			imported++
		}
	}

	// Summary
	fmt.Println("\n✅ Import complete!")
	fmt.Printf("   Imported: %d\n", imported)
	if skipped > 0 {
		fmt.Printf("   Skipped (duplicates): %d\n", skipped)
	}
	if errors > 0 {
		fmt.Printf("   Errors: %d\n", errors)
	}

	if autoCategorize {
		categorized := 0
		for _, tx := range transactions {
			if tx.Category != "Uncategorized" && tx.Category != "" {
				categorized++
			}
		}
		fmt.Printf("   Auto-categorized: %d/%d\n", categorized, imported)
	}
}
