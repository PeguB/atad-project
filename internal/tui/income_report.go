package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/PeguB/atad-project/internal/repository"
	tea "github.com/charmbracelet/bubbletea"
)

type IncomeReportScreen struct {
	repo         *repository.TransactionRepository
	totalIncome  float64
	byCategory   map[string]float64
	startDate    string
	endDate      string
	step         int // 0: select period, 1: custom dates, 2: show results
	err          string
	periodChoice int // 0: all time, 1: this month, 2: this year, 3: custom
}

func NewIncomeReportScreen(repo *repository.TransactionRepository) *IncomeReportScreen {
	return &IncomeReportScreen{
		repo: repo,
		step: 0,
	}
}

func (s *IncomeReportScreen) Init() {
	s.totalIncome = 0
	s.byCategory = make(map[string]float64)
	s.startDate = ""
	s.endDate = ""
	s.step = 0
	s.err = ""
	s.periodChoice = 0
}

func (s *IncomeReportScreen) Update(msg tea.Msg) (*IncomeReportScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.step {
		case 0: // Select period
			switch msg.String() {
			case "1": // All time
				s.periodChoice = 0
				s.calculateIncome()
				s.step = 2
			case "2": // This month
				s.periodChoice = 1
				now := time.Now()
				s.startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("02/01/2006")
				endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
				s.endDate = endOfMonth.Format("02/01/2006")
				s.calculateIncome()
				s.step = 2
			case "3": // This year
				s.periodChoice = 2
				now := time.Now()
				s.startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("02/01/2006")
				s.endDate = time.Date(now.Year(), 12, 31, 0, 0, 0, 0, now.Location()).Format("02/01/2006")
				s.calculateIncome()
				s.step = 2
			case "4": // Custom
				s.periodChoice = 3
				s.step = 1
			}
		case 1: // Enter custom date (for now, skip to results - can be enhanced)
			switch msg.String() {
			case "enter":
				// For now, default to this month if not implemented
				now := time.Now()
				s.startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("02/01/2006")
				endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
				s.endDate = endOfMonth.Format("02/01/2006")
				s.calculateIncome()
				s.step = 2
			}
		case 2: // Show results
			// ESC handled by parent
		}
	}
	return s, nil
}

func (s *IncomeReportScreen) calculateIncome() {
	s.totalIncome = 0
	s.byCategory = make(map[string]float64)
	s.err = ""

	// Get all transactions
	transactions, err := s.repo.GetAll()
	if err != nil {
		s.err = fmt.Sprintf("Error loading transactions: %v", err)
		return
	}

	// Parse date filters if set
	var startTime, endTime time.Time
	var useFilter bool
	if s.startDate != "" && s.endDate != "" {
		startTime, err = time.Parse("02/01/2006", s.startDate)
		if err != nil {
			s.err = fmt.Sprintf("Invalid start date: %v", err)
			return
		}
		endTime, err = time.Parse("02/01/2006", s.endDate)
		if err != nil {
			s.err = fmt.Sprintf("Invalid end date: %v", err)
			return
		}
		useFilter = true
	}

	// Calculate totals
	for _, tx := range transactions {
		if tx.Type != "income" {
			continue
		}

		// Apply date filter if set
		if useFilter {
			if tx.Date.Before(startTime) || tx.Date.After(endTime) {
				continue
			}
		}

		s.totalIncome += tx.Amount
		s.byCategory[tx.Category] += tx.Amount
	}
}

func (s *IncomeReportScreen) View() string {
	var b strings.Builder

	b.WriteString("📊 Income Report\n\n")

	switch s.step {
	case 0:
		b.WriteString("Select time period:\n\n")
		b.WriteString("  1. All time\n")
		b.WriteString("  2. This month\n")
		b.WriteString("  3. This year\n")
		b.WriteString("  4. Custom period\n")
		b.WriteString("\nPress ESC to return to menu\n")
	case 1:
		b.WriteString("Custom period selected\n\n")
		b.WriteString("(Press Enter to use current month - custom dates coming soon)\n")
		b.WriteString("\nPress ESC to cancel\n")
	case 2:
		if s.err != "" {
			b.WriteString(fmt.Sprintf("❌ %s\n\n", s.err))
		} else {
			// Show period
			periodName := ""
			switch s.periodChoice {
			case 0:
				periodName = "All Time"
			case 1:
				periodName = "This Month"
			case 2:
				periodName = "This Year"
			case 3:
				periodName = "Custom Period"
			}

			if s.startDate != "" && s.endDate != "" {
				b.WriteString(fmt.Sprintf("Period: %s (%s - %s)\n\n", periodName, s.startDate, s.endDate))
			} else {
				b.WriteString(fmt.Sprintf("Period: %s\n\n", periodName))
			}

			// Show total
			b.WriteString(fmt.Sprintf("💰 Total Income: $%.2f\n\n", s.totalIncome))

			// Show breakdown by category
			if len(s.byCategory) > 0 {
				b.WriteString("Breakdown by Category:\n")
				b.WriteString("─────────────────────────────────────\n")

				for category, amount := range s.byCategory {
					percentage := (amount / s.totalIncome) * 100
					b.WriteString(fmt.Sprintf("  %-20s $%8.2f  (%.1f%%)\n", category, amount, percentage))
				}
			} else {
				b.WriteString("No income transactions found for this period.\n")
			}
		}

		b.WriteString("\nPress ESC to return to menu\n")
	}

	return b.String()
}

func (s *IncomeReportScreen) Reset() {
	s.Init()
}
