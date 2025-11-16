package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PeguB/atad-project/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	db       *database.Database
	choices  []string
	cursor   int
	selected map[int]struct{}
	status   string
}

func initialModel() model {
	return model{
		choices:  []string{"Test Database Connection", "View Transactions", "Add Transaction", "Exit"},
		selected: make(map[int]struct{}),
		status:   "Ready",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					if err := db.DB.Ping(); err != nil {
						m.status = fmt.Sprintf("❌ Connection failed: %v", err)
					} else {
						m.status = "✅ Database connection successful!"
					}
				}
			case 1: // View Transactions
				m.status = "📋 View Transactions - Coming soon!"
			case 2: // Add Transaction
				m.status = "➕ Add Transaction - Coming soon!"
			case 3: // Exit
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
