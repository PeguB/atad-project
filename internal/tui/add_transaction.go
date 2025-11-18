package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PeguB/atad-project/internal/models"
	"github.com/PeguB/atad-project/internal/repository"
	"github.com/PeguB/atad-project/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

type AddTransactionScreen struct {
	repo            *repository.TransactionRepository
	categoryService *service.CategoryService
	step            int
	txType          string
	description     string
	amount          string
	category        string
	date            string // Date in DD/MM/YYYY format
	suggestedCat    string
	err             string
	success         string
}

func NewAddTransactionScreen(repo *repository.TransactionRepository, categoryService *service.CategoryService) *AddTransactionScreen {
	return &AddTransactionScreen{
		repo:            repo,
		categoryService: categoryService,
		step:            0,
	}
}

func (s *AddTransactionScreen) Update(msg tea.Msg) (*AddTransactionScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.step {
		case 0: // Select type
			switch msg.String() {
			case "1":
				s.txType = "income"
				s.step = 1
			case "2":
				s.txType = "expense"
				s.step = 1
			}
		case 1: // Enter description
			switch msg.String() {
			case "enter":
				if s.description != "" {
					s.step = 2
				}
			case "backspace":
				if len(s.description) > 0 {
					s.description = s.description[:len(s.description)-1]
				}
			default:
				if len(msg.String()) == 1 {
					s.description += msg.String()
				}
			}
		case 2: // Enter amount
			switch msg.String() {
			case "enter":
				if _, err := strconv.ParseFloat(s.amount, 64); err == nil && s.amount != "" {
					s.step = 3
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
		case 3: // Enter date
			switch msg.String() {
			case "enter":
				if s.date == "" {
					// Use today's date if nothing entered
					s.date = time.Now().Format("02/01/2006")
				}
				if _, err := time.Parse("02/01/2006", s.date); err == nil {
					// Auto-categorize based on description
					s.suggestedCat = s.categoryService.CategorizeTransaction(s.description)
					s.step = 4
					s.err = ""
				} else {
					s.err = "Invalid date format (use DD/MM/YYYY)"
				}
			case "backspace":
				if len(s.date) > 0 {
					s.date = s.date[:len(s.date)-1]
				}
			default:
				if len(msg.String()) == 1 && (msg.String()[0] >= '0' && msg.String()[0] <= '9' || msg.String() == "/") {
					if len(s.date) < 10 {
						s.date += msg.String()
					}
				}
			}
		case 4: // Enter category
			switch msg.String() {
			case "enter":
				// If user didn't type anything, use suggested category
				if s.category == "" && s.suggestedCat != "" {
					s.category = s.suggestedCat
				}
				if s.category != "" {
					s.saveTransaction()
				}
			case "tab":
				// Accept suggestion with tab key
				if s.suggestedCat != "" {
					s.category = s.suggestedCat
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
		}
	}
	return s, nil
}

func (s *AddTransactionScreen) saveTransaction() {
	amount, _ := strconv.ParseFloat(s.amount, 64)

	// Parse the date
	txDate, err := time.Parse("02/01/2006", s.date)
	if err != nil {
		s.err = fmt.Sprintf("Invalid date: %v", err)
		s.step = 5
		return
	}

	tx := &models.Transaction{
		Date:        txDate,
		Description: s.description,
		Amount:      amount,
		Category:    s.category,
		Type:        s.txType,
	}

	err = s.repo.Create(tx)
	if err != nil {
		s.err = fmt.Sprintf("Failed to save: %v", err)
		s.step = 5
	} else {
		s.success = fmt.Sprintf("✅ %s added: $%.2f (%s) on %s",
			strings.Title(s.txType), amount, s.category, s.date)
		s.step = 5
	}
}

func (s *AddTransactionScreen) View() string {
	var b strings.Builder

	b.WriteString("➕ Add Transaction\n\n")

	switch s.step {
	case 0:
		b.WriteString("Select transaction type:\n\n")
		b.WriteString("  1. Income\n")
		b.WriteString("  2. Expense\n")
	case 1:
		b.WriteString(fmt.Sprintf("Type: %s\n\n", strings.Title(s.txType)))
		b.WriteString("Description: " + s.description + "▊\n")
		b.WriteString("\n(Press Enter to continue)\n")
	case 2:
		b.WriteString(fmt.Sprintf("Type: %s\n", strings.Title(s.txType)))
		b.WriteString(fmt.Sprintf("Description: %s\n\n", s.description))
		b.WriteString("Amount: $" + s.amount + "▊\n")
		if s.err != "" {
			b.WriteString("\n❌ " + s.err + "\n")
		}
		b.WriteString("\n(Press Enter to continue)\n")
	case 3:
		b.WriteString(fmt.Sprintf("Type: %s\n", strings.Title(s.txType)))
		b.WriteString(fmt.Sprintf("Description: %s\n", s.description))
		b.WriteString(fmt.Sprintf("Amount: $%s\n\n", s.amount))
		if s.date == "" {
			b.WriteString("Date (DD/MM/YYYY): ▊\n")
			b.WriteString("\n(Press Enter for today's date)\n")
		} else {
			b.WriteString("Date (DD/MM/YYYY): " + s.date + "▊\n")
			if s.err != "" {
				b.WriteString("\n❌ " + s.err + "\n")
			}
			b.WriteString("\n(Press Enter to continue)\n")
		}
	case 4:
		b.WriteString(fmt.Sprintf("Type: %s\n", strings.Title(s.txType)))
		b.WriteString(fmt.Sprintf("Description: %s\n", s.description))
		b.WriteString(fmt.Sprintf("Amount: $%s\n", s.amount))
		b.WriteString(fmt.Sprintf("Date: %s\n\n", s.date))

		if s.category == "" && s.suggestedCat != "" {
			b.WriteString(fmt.Sprintf("💡 Suggested: %s\n", s.suggestedCat))
		}

		b.WriteString("Category: " + s.category + "▊\n")
		b.WriteString("\n(Press Enter to save")
		if s.category == "" && s.suggestedCat != "" {
			b.WriteString(" with suggestion")
		}
		b.WriteString(", Tab to accept suggestion)\n")
	case 5:
		if s.success != "" {
			b.WriteString(s.success + "\n\n")
		}
		if s.err != "" {
			b.WriteString("❌ " + s.err + "\n\n")
		}
		b.WriteString("Press ESC to return to menu\n")
	}

	return b.String()
}

func (s *AddTransactionScreen) Reset() {
	s.step = 0
	s.txType = ""
	s.description = ""
	s.amount = ""
	s.category = ""
	s.date = ""
	s.suggestedCat = ""
	s.err = ""
	s.success = ""
}
