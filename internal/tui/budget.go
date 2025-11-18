package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PeguB/atad-project/internal/models"
	"github.com/PeguB/atad-project/internal/repository"
	tea "github.com/charmbracelet/bubbletea"
)

type BudgetScreen struct {
	budgetRepo *repository.BudgetRepository
	mode       string // "list", "add", "view"
	budgets    []*models.Budget

	// Add budget fields
	step      int
	category  string
	amount    string
	startDate string // Date in DD/MM/YYYY format
	endDate   string // Date in DD/MM/YYYY format

	err     string
	success string
}

func NewBudgetScreen(budgetRepo *repository.BudgetRepository) *BudgetScreen {
	return &BudgetScreen{
		budgetRepo: budgetRepo,
		mode:       "list",
	}
}

func (s *BudgetScreen) Init() error {
	budgets, err := s.budgetRepo.GetAll()
	if err != nil {
		return err
	}
	s.budgets = budgets
	return nil
}

func (s *BudgetScreen) Update(msg tea.Msg) (*BudgetScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.mode == "list" {
			switch msg.String() {
			case "n":
				s.mode = "add"
				s.step = 0
				s.resetAddFields()
			}
		} else if s.mode == "add" {
			return s.handleAddBudget(msg)
		}
	}
	return s, nil
}

func (s *BudgetScreen) handleAddBudget(msg tea.KeyMsg) (*BudgetScreen, tea.Cmd) {
	switch s.step {
	case 0: // Enter category
		switch msg.String() {
		case "enter":
			if s.category != "" {
				s.step = 1
			}
		case "backspace":
			if len(s.category) > 0 {
				s.category = s.category[:len(s.category)-1]
			}
		default:
			if len(msg.String()) == 1 {
				s.category += msg.String()
			}
		}
	case 1: // Enter amount
		switch msg.String() {
		case "enter":
			if _, err := strconv.ParseFloat(s.amount, 64); err == nil && s.amount != "" {
				s.step = 2
				s.err = ""
			} else {
				s.err = "Invalid amount"
			}
		case "backspace":
			if len(s.amount) > 0 {
				s.amount = s.amount[:len(s.amount)-1]
			}
		default:
			if len(msg.String()) == 1 && (msg.String()[0] >= '0' && msg.String()[0] <= '9' || msg.String() == ".") {
				s.amount += msg.String()
			}
		}
	case 2: // Enter start date
		switch msg.String() {
		case "enter":
			if s.startDate != "" {
				if _, err := time.Parse("02/01/2006", s.startDate); err == nil {
					s.step = 3
					s.err = ""
				} else {
					s.err = "Invalid date (e.g., Sept has 30 days, not 31)"
				}
			} else {
				s.err = "Date required (use DD/MM/YYYY)"
			}
		case "backspace":
			if len(s.startDate) > 0 {
				s.startDate = s.startDate[:len(s.startDate)-1]
				s.err = ""
			}
		default:
			if len(msg.String()) == 1 && (msg.String()[0] >= '0' && msg.String()[0] <= '9' || msg.String() == "/") {
				if len(s.startDate) < 10 {
					s.startDate += msg.String()
					s.err = ""
				}
			}
		}
	case 3: // Enter end date
		switch msg.String() {
		case "enter":
			if s.endDate != "" {
				if _, err := time.Parse("02/01/2006", s.endDate); err == nil {
					s.err = ""
					s.saveBudget()
				} else {
					s.err = "Invalid date (e.g., Sept has 30 days, not 31)"
				}
			} else {
				s.err = "Date required (use DD/MM/YYYY)"
			}
		case "backspace":
			if len(s.endDate) > 0 {
				s.endDate = s.endDate[:len(s.endDate)-1]
				s.err = ""
			}
		default:
			if len(msg.String()) == 1 && (msg.String()[0] >= '0' && msg.String()[0] <= '9' || msg.String() == "/") {
				if len(s.endDate) < 10 {
					s.endDate += msg.String()
					s.err = ""
				}
			}
		}
	}
	return s, nil
}

func (s *BudgetScreen) saveBudget() {
	amount, _ := strconv.ParseFloat(s.amount, 64)

	startDate, err := time.Parse("02/01/2006", s.startDate)
	if err != nil {
		s.err = fmt.Sprintf("Invalid start date: %v", err)
		s.step = 2 // Go back to start date step
		return
	}
	endDate, err := time.Parse("02/01/2006", s.endDate)
	if err != nil {
		s.err = fmt.Sprintf("Invalid end date: %v", err)
		s.step = 3 // Stay on end date step
		return
	}

	budget := &models.Budget{
		Category:  s.category,
		Amount:    amount,
		Period:    "custom",
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Check if budget already exists for this category and date range
	existing, err := s.budgetRepo.GetByCategoryAndDateRange(s.category, startDate, endDate)
	if err != nil {
		s.err = fmt.Sprintf("Failed to check existing budget: %v", err)
		s.mode = "list"
		s.resetAddFields()
		return
	}

	if existing != nil {
		// Update existing budget
		err = s.budgetRepo.Update(budget)
		if err != nil {
			s.err = fmt.Sprintf("Failed to update: %v", err)
		} else {
			s.success = fmt.Sprintf("✅ Budget updated successfully! $%.2f for %s (%s to %s)", amount, s.category, s.startDate, s.endDate)
			s.Init() // Reload budgets
		}
	} else {
		// Create new budget
		err = s.budgetRepo.Create(budget)
		if err != nil {
			s.err = fmt.Sprintf("Failed to save: %v", err)
		} else {
			s.success = fmt.Sprintf("✅ Budget created successfully! $%.2f for %s (%s to %s)", amount, s.category, s.startDate, s.endDate)
			s.Init() // Reload budgets
		}
	}
	s.Init()
	s.mode = "list"
	s.resetAddFields()
}

func (s *BudgetScreen) resetAddFields() {
	s.step = 0
	s.category = ""
	s.amount = ""
	s.startDate = ""
	s.endDate = ""
	s.err = ""
}

func (s *BudgetScreen) View() string {
	var b strings.Builder

	if s.mode == "list" {
		b.WriteString("💰 Budget Management\n\n")

		if len(s.budgets) == 0 {
			b.WriteString("No budgets set yet.\n\n")
		} else {
			b.WriteString(fmt.Sprintf("%-20s %-15s %-12s %s\n", "Category", "Period", "Budget", "Status"))
			b.WriteString(strings.Repeat("-", 75) + "\n")

			for _, budget := range s.budgets {
				// Get current spending
				spending, err := s.budgetRepo.GetSpending(budget.Category, budget.StartDate, budget.EndDate)
				percentage := 0.0
				if budget.Amount > 0 {
					percentage = (spending / budget.Amount) * 100
				}

				status := "✅"
				if percentage > 100 {
					status = "🚨"
				} else if percentage > 80 {
					status = "⚠️"
				}

				spendingInfo := ""
				if err == nil {
					spendingInfo = fmt.Sprintf(" %s $%.2f/%.2f (%.0f%%)", status, spending, budget.Amount, percentage)
				}

				periodStr := fmt.Sprintf("%s-%s", budget.StartDate.Format("02/01/06"), budget.EndDate.Format("02/01/06"))

				b.WriteString(fmt.Sprintf("%-20s %-15s $%-11.2f%s\n",
					budget.Category, periodStr, budget.Amount, spendingInfo))
			}
			b.WriteString("\n")
		}

		if s.success != "" {
			b.WriteString(s.success + "\n\n")
			s.success = ""
		}

		b.WriteString("\nPress 'n' to add new budget | ESC to return\n")

	} else if s.mode == "add" {
		b.WriteString("➕ Add Budget\n\n")

		switch s.step {
		case 0:
			b.WriteString("Category: " + s.category + "▊\n")
			b.WriteString("\n(Press Enter to continue)\n")
		case 1:
			b.WriteString(fmt.Sprintf("Category: %s\n\n", s.category))
			b.WriteString("Budget Amount: $" + s.amount + "▊\n")
			if s.err != "" {
				b.WriteString("\n❌ " + s.err + "\n")
			}
			b.WriteString("\n(Press Enter to continue)\n")
		case 2:
			b.WriteString(fmt.Sprintf("Category: %s\n", s.category))
			b.WriteString(fmt.Sprintf("Amount: $%s\n\n", s.amount))
			b.WriteString("Start Date (DD/MM/YYYY): " + s.startDate + "▊\n")
			if s.err != "" {
				b.WriteString("\n❌ " + s.err + "\n")
			}
			b.WriteString("\n(Press Enter to continue)\n")
		case 3:
			b.WriteString(fmt.Sprintf("Category: %s\n", s.category))
			b.WriteString(fmt.Sprintf("Amount: $%s\n", s.amount))
			b.WriteString(fmt.Sprintf("Start Date: %s\n\n", s.startDate))
			b.WriteString("End Date (DD/MM/YYYY): " + s.endDate + "▊\n")
			if s.err != "" {
				b.WriteString("\n❌ " + s.err + "\n")
			}
			b.WriteString("\n(Press Enter to save)\n")
		}

		if s.err != "" && s.step != 1 {
			b.WriteString("\n❌ " + s.err + "\n")
		}
	}

	return b.String()
}

func (s *BudgetScreen) Reset() {
	s.mode = "list"
	s.resetAddFields()
	s.Init()
}
