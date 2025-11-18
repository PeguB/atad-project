package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PeguB/atad-project/internal/database"
	"github.com/PeguB/atad-project/internal/repository"
	"github.com/PeguB/atad-project/internal/service"
	"github.com/PeguB/atad-project/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	menuScreen screen = iota
	viewTransactionsScreen
	addTransactionScreen
	budgetScreen
)

type model struct {
	db                     *database.Database
	repo                   *repository.TransactionRepository
	budgetRepo             *repository.BudgetRepository
	categoryService        *service.CategoryService
	currentScreen          screen
	viewTransactionsScreen *tui.ViewTransactionsScreen
	addTransactionScreen   *tui.AddTransactionScreen
	budgetScreen           *tui.BudgetScreen
	choices                []string
	cursor                 int
	selected               map[int]struct{}
	status                 string
}

func initialModel() model {
	return model{
		currentScreen:   menuScreen,
		categoryService: service.NewCategoryService(),
		choices:         []string{"Test Database Connection", "View Transactions", "Add Transaction", "Manage Budgets", "Exit"},
		selected:        make(map[int]struct{}),
		status:          "Ready",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle screen switching
	if m.currentScreen == viewTransactionsScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.viewTransactionsScreen.Reset()
				m.currentScreen = menuScreen
				m.status = "Returned to menu"
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.viewTransactionsScreen, cmd = m.viewTransactionsScreen.Update(msg)
		return m, cmd
	}

	if m.currentScreen == addTransactionScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.addTransactionScreen.Reset()
				m.currentScreen = menuScreen
				m.status = "Transaction cancelled"
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.addTransactionScreen, cmd = m.addTransactionScreen.Update(msg)
		return m, cmd
	}

	if m.currentScreen == budgetScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.budgetScreen.Reset()
				m.currentScreen = menuScreen
				m.status = "Returned to menu"
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.budgetScreen, cmd = m.budgetScreen.Update(msg)
		return m, cmd
	}

	// Main menu handling
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.db != nil {
				m.db.Close()
			}
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			switch m.cursor {
			case 0: // Test Database Connection
				db, err := database.NewDatabase()
				if err != nil {
					m.status = fmt.Sprintf("❌ Failed to connect: %v", err)
				} else {
					m.db = db
					m.repo = repository.NewTransactionRepository(db.DB)
					m.budgetRepo = repository.NewBudgetRepository(db.DB)
					if err := db.DB.Ping(); err != nil {
						m.status = fmt.Sprintf("❌ Connection failed: %v", err)
					} else {
						m.status = "✅ Database connection successful!"
					}
				}
			case 1: // View Transactions
				if m.db == nil {
					db, err := database.NewDatabase()
					if err != nil {
						m.status = fmt.Sprintf("❌ Failed to connect to database: %v", err)
						return m, nil
					}
					m.db = db
					m.repo = repository.NewTransactionRepository(db.DB)
					m.budgetRepo = repository.NewBudgetRepository(db.DB)
				}
				m.viewTransactionsScreen = tui.NewViewTransactionsScreen(m.repo, m.budgetRepo)
				m.viewTransactionsScreen.Init()
				m.currentScreen = viewTransactionsScreen
			case 2: // Add Transaction
				if m.db == nil {
					db, err := database.NewDatabase()
					if err != nil {
						m.status = fmt.Sprintf("❌ Failed to connect to database: %v", err)
						return m, nil
					}
					m.db = db
					m.repo = repository.NewTransactionRepository(db.DB)
					m.budgetRepo = repository.NewBudgetRepository(db.DB)
				}
				m.addTransactionScreen = tui.NewAddTransactionScreen(m.repo, m.categoryService)
				m.currentScreen = addTransactionScreen
			case 3: // Manage Budgets
				if m.db == nil {
					db, err := database.NewDatabase()
					if err != nil {
						m.status = fmt.Sprintf("❌ Failed to connect to database: %v", err)
						return m, nil
					}
					m.db = db
					m.repo = repository.NewTransactionRepository(db.DB)
					m.budgetRepo = repository.NewBudgetRepository(db.DB)
				}
				m.budgetScreen = tui.NewBudgetScreen(m.budgetRepo)
				m.budgetScreen.Init()
				m.currentScreen = budgetScreen
			case 4: // Exit
				if m.db != nil {
					m.db.Close()
				}
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.currentScreen == viewTransactionsScreen {
		statusMsg := ""
		if m.status != "" && m.status != "Ready" {
			statusMsg = fmt.Sprintf("\nStatus: %s\n", m.status)
		}
		return m.viewTransactionsScreen.View() + statusMsg
	}

	if m.currentScreen == addTransactionScreen {
		statusMsg := ""
		if m.status != "" && m.status != "Ready" {
			statusMsg = fmt.Sprintf("\nStatus: %s", m.status)
		}
		return m.addTransactionScreen.View() + "\n\nPress ESC to return to menu" + statusMsg
	}

	if m.currentScreen == budgetScreen {
		statusMsg := ""
		if m.status != "" && m.status != "Ready" {
			statusMsg = fmt.Sprintf("\nStatus: %s\n", m.status)
		}
		return m.budgetScreen.View() + statusMsg
	}

	s := "🏦 ATAD - Personal Finance Tracker\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += fmt.Sprintf("\nStatus: %s\n", m.status)
	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
