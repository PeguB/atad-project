package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PeguB/atad-project/internal/models"
	"github.com/PeguB/atad-project/internal/repository"
	tea "github.com/charmbracelet/bubbletea"
)

type sortField int

const (
	sortByDate sortField = iota
	sortByAmount
	sortByCategory
	sortByDescription
)

type ViewTransactionsScreen struct {
	repo         *repository.TransactionRepository
	budgetRepo   *repository.BudgetRepository
	transactions []*models.Transaction
	budgets      []*models.Budget
	cursor       int
	page         int
	pageSize     int
	sortBy       sortField
	sortDesc     bool
	filterType   string // "all", "income", "expense"
	filterText   string
	filterMode   bool // true when entering filter text
	showBudgets  bool // true when 'b' is pressed to show budgets
	err          error
	total        float64
}

func NewViewTransactionsScreen(repo *repository.TransactionRepository, budgetRepo *repository.BudgetRepository) *ViewTransactionsScreen {
	return &ViewTransactionsScreen{
		repo:       repo,
		budgetRepo: budgetRepo,
		cursor:     0,
		page:       0,
		pageSize:   15,
		sortBy:     sortByDate,
		sortDesc:   true,
		filterType: "all",
	}
}

func (s *ViewTransactionsScreen) Init() {
	s.loadTransactions()
	s.loadBudgets()
}

func (s *ViewTransactionsScreen) loadTransactions() {
	var err error

	if s.filterType == "all" {
		s.transactions, err = s.repo.GetAll()
	} else {
		s.transactions, err = s.repo.GetByType(s.filterType)
	}

	if err != nil {
		s.err = err
		return
	}

	// Apply text filter
	if s.filterText != "" {
		filtered := []*models.Transaction{}
		searchLower := strings.ToLower(s.filterText)
		for _, tx := range s.transactions {
			if strings.Contains(strings.ToLower(tx.Description), searchLower) ||
				strings.Contains(strings.ToLower(tx.Category), searchLower) {
				filtered = append(filtered, tx)
			}
		}
		s.transactions = filtered
	}

	// Sort transactions
	s.sortTransactions()

	// Calculate total
	s.total = 0
	for _, tx := range s.transactions {
		if tx.Type == "income" {
			s.total += tx.Amount
		} else {
			s.total -= tx.Amount
		}
	}

	// Reset cursor if needed
	if s.cursor >= len(s.transactions) {
		s.cursor = 0
	}

	s.err = nil
}

func (s *ViewTransactionsScreen) loadBudgets() {
	if s.budgetRepo == nil {
		return
	}
	budgets, err := s.budgetRepo.GetAll()
	if err != nil {
		// Silently ignore budget loading errors
		return
	}
	s.budgets = budgets
}

func (s *ViewTransactionsScreen) sortTransactions() {
	sort.Slice(s.transactions, func(i, j int) bool {
		var less bool
		switch s.sortBy {
		case sortByDate:
			less = s.transactions[i].Date.Before(s.transactions[j].Date)
		case sortByAmount:
			less = s.transactions[i].Amount < s.transactions[j].Amount
		case sortByCategory:
			less = s.transactions[i].Category < s.transactions[j].Category
		case sortByDescription:
			less = s.transactions[i].Description < s.transactions[j].Description
		}

		if s.sortDesc {
			return !less
		}
		return less
	})
}

func (s *ViewTransactionsScreen) Update(msg tea.Msg) (*ViewTransactionsScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Filter mode handling
		if s.filterMode {
			switch msg.String() {
			case "enter", "esc":
				s.filterMode = false
				s.loadTransactions()
			case "backspace":
				if len(s.filterText) > 0 {
					s.filterText = s.filterText[:len(s.filterText)-1]
					s.loadTransactions()
				}
			default:
				if len(msg.String()) == 1 {
					s.filterText += msg.String()
					s.loadTransactions()
				}
			}
			return s, nil
		}

		// Normal mode handling
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
				if s.cursor < s.page*s.pageSize {
					s.page--
				}
			}
		case "down", "j":
			if s.cursor < len(s.transactions)-1 {
				s.cursor++
				if s.cursor >= (s.page+1)*s.pageSize {
					s.page++
				}
			}
		case "g": // Go to top
			s.cursor = 0
			s.page = 0
		case "G": // Go to bottom
			if len(s.transactions) > 0 {
				s.cursor = len(s.transactions) - 1
				s.page = s.cursor / s.pageSize
			}
		case "d": // Sort by date
			if s.sortBy == sortByDate {
				s.sortDesc = !s.sortDesc
			} else {
				s.sortBy = sortByDate
				s.sortDesc = true
			}
			s.loadTransactions()
		case "a": // Sort by amount
			if s.sortBy == sortByAmount {
				s.sortDesc = !s.sortDesc
			} else {
				s.sortBy = sortByAmount
				s.sortDesc = true
			}
			s.loadTransactions()
		case "c": // Sort by category
			if s.sortBy == sortByCategory {
				s.sortDesc = !s.sortDesc
			} else {
				s.sortBy = sortByCategory
				s.sortDesc = false
			}
			s.loadTransactions()
		case "n": // Sort by description
			if s.sortBy == sortByDescription {
				s.sortDesc = !s.sortDesc
			} else {
				s.sortBy = sortByDescription
				s.sortDesc = false
			}
			s.loadTransactions()
		case "f": // Filter by type
			switch s.filterType {
			case "all":
				s.filterType = "income"
			case "income":
				s.filterType = "expense"
			case "expense":
				s.filterType = "all"
			}
			s.loadTransactions()
		case "s": // Search/filter text
			s.filterMode = true
		case "x": // Clear filters
			s.filterText = ""
			s.filterType = "all"
			s.loadTransactions()
		case "b": // Toggle budgets view
			s.showBudgets = !s.showBudgets
		case "r": // Refresh
			s.loadTransactions()
			s.loadBudgets()
		case "delete", "backspace":
			if len(s.transactions) > 0 && s.cursor < len(s.transactions) {
				tx := s.transactions[s.cursor]
				err := s.repo.Delete(tx.ID)
				if err == nil {
					s.loadTransactions()
				} else {
					s.err = err
				}
			}
		}
	}

	return s, nil
}

func (s *ViewTransactionsScreen) View() string {
	var b strings.Builder

	b.WriteString("📋 View Transactions\n\n")

	if s.err != nil {
		b.WriteString(fmt.Sprintf("❌ Error: %v\n", s.err))
		return b.String()
	}

	// Show budgets if toggled on
	if s.showBudgets {
		b.WriteString("💰 Active Budgets\n")
		b.WriteString("  ────────────────────────────────────────────────────────────────────────────\n")
		if len(s.budgets) == 0 {
			b.WriteString("  No budgets found.\n")
		} else {
			b.WriteString("  Category             Budget      Spent       Remaining   Period\n")
			b.WriteString("  ──────────────────── ─────────── ─────────── ─────────── ──────────────────\n")
			for _, budget := range s.budgets {
				cat := budget.Category
				if len(cat) > 20 {
					cat = cat[:17] + "..."
				}

				// Get spending for this budget
				spent, _ := s.budgetRepo.GetSpending(budget.Category, budget.StartDate, budget.EndDate)
				remaining := budget.Amount - spent

				// Format amounts
				budgetStr := fmt.Sprintf("$%.2f", budget.Amount)
				spentStr := fmt.Sprintf("$%.2f", spent)
				remainingStr := fmt.Sprintf("$%.2f", remaining)

				// Format period
				periodStr := fmt.Sprintf("%s - %s",
					budget.StartDate.Format("02/01/2006"),
					budget.EndDate.Format("02/01/2006"))

				b.WriteString(fmt.Sprintf("  %-20s %-11s %-11s %-11s %s\n",
					cat, budgetStr, spentStr, remainingStr, periodStr))
			}
		}
		b.WriteString("  ────────────────────────────────────────────────────────────────────────────\n\n")
	}

	if len(s.transactions) == 0 {
		b.WriteString("No transactions found.\n")
		b.WriteString("\nPress ESC to return to menu")
		return b.String()
	}

	// Display filters and sort info
	b.WriteString(fmt.Sprintf("Filter: %s", s.filterType))
	if s.filterText != "" {
		b.WriteString(fmt.Sprintf(" | Search: \"%s\"", s.filterText))
	}
	b.WriteString(" | Sort: ")
	switch s.sortBy {
	case sortByDate:
		b.WriteString("Date")
	case sortByAmount:
		b.WriteString("Amount")
	case sortByCategory:
		b.WriteString("Category")
	case sortByDescription:
		b.WriteString("Description")
	}
	if s.sortDesc {
		b.WriteString(" ↓")
	} else {
		b.WriteString(" ↑")
	}
	b.WriteString("\n\n")

	// Column headers
	b.WriteString("  Date       Type     Amount      Category             Description\n")
	b.WriteString("  ────────── ──────── ─────────── ──────────────────── ────────────────────────────\n")

	// Display transactions for current page
	start := s.page * s.pageSize
	end := start + s.pageSize
	if end > len(s.transactions) {
		end = len(s.transactions)
	}

	for i := start; i < end; i++ {
		tx := s.transactions[i]
		cursor := " "
		if i == s.cursor {
			cursor = ">"
		}

		// Format amount with sign
		amountStr := fmt.Sprintf("$%.2f", tx.Amount)
		if tx.Type == "income" {
			amountStr = "+" + amountStr
		} else {
			amountStr = "-" + amountStr
		}

		// Format type
		typeStr := tx.Type
		if tx.Type == "income" {
			typeStr = "Income"
		} else {
			typeStr = "Expense"
		}

		// Truncate description if too long
		desc := tx.Description
		if len(desc) > 30 {
			desc = desc[:27] + "..."
		}

		// Truncate category if too long
		cat := tx.Category
		if len(cat) > 20 {
			cat = cat[:17] + "..."
		}

		b.WriteString(fmt.Sprintf("%s %s %-8s %-11s %-20s %s\n",
			cursor,
			tx.Date.Format("2006-01-02"),
			typeStr,
			amountStr,
			cat,
			desc))
	}

	// Summary line
	b.WriteString("  ────────────────────────────────────────────────────────────────────────────\n")
	totalStr := fmt.Sprintf("$%.2f", s.total)
	if s.total >= 0 {
		totalStr = "+" + totalStr
	}
	b.WriteString(fmt.Sprintf("  Total (%d transactions): %s\n", len(s.transactions), totalStr))

	// Pagination info
	if len(s.transactions) > s.pageSize {
		totalPages := (len(s.transactions) + s.pageSize - 1) / s.pageSize
		b.WriteString(fmt.Sprintf("\n  Page %d/%d (showing %d-%d of %d)\n",
			s.page+1, totalPages, start+1, end, len(s.transactions)))
	}

	// Help text
	b.WriteString("\n")
	if s.filterMode {
		b.WriteString("🔍 Search mode - Type to search, Enter to confirm, ESC to cancel\n")
		b.WriteString(fmt.Sprintf("Current search: %s_\n", s.filterText))
	} else {
		b.WriteString("Navigation: ↑/↓ or k/j | g/G = top/bottom | d/a/c/n = sort by Date/Amount/Category/Name\n")
		b.WriteString("Filter: f = cycle type | s = search | x = clear filters | b = toggle budgets | r = refresh | Delete = remove\n")
		b.WriteString("Press ESC to return to menu | q to quit\n")
	}

	return b.String()
}

func (s *ViewTransactionsScreen) Reset() {
	s.cursor = 0
	s.page = 0
	s.filterText = ""
	s.filterType = "all"
	s.filterMode = false
	s.sortBy = sortByDate
	s.sortDesc = true
}
