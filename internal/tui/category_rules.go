package tui

import (
	"fmt"
	"strings"

	"github.com/PeguB/atad-project/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

type CategoryRulesScreen struct {
	categoryService *service.CategoryService
}

func NewCategoryRulesScreen(categoryService *service.CategoryService) *CategoryRulesScreen {
	return &CategoryRulesScreen{
		categoryService: categoryService,
	}
}

func (s *CategoryRulesScreen) Update(msg tea.Msg) (*CategoryRulesScreen, tea.Cmd) {
	return s, nil
}

func (s *CategoryRulesScreen) View() string {
	var b strings.Builder

	b.WriteString("🏷️  Automatic Categorization Rules\n\n")
	b.WriteString("These patterns are used to automatically categorize your transactions:\n\n")

	rules := s.categoryService.GetRules()

	for i, rule := range rules {
		if i > 0 && i%5 == 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("• %-15s - %s\n", rule.Category, rule.Description))
	}

	b.WriteString("\n\nExample: Typing 'Starbucks Coffee' will automatically suggest 'Restaurants'\n")
	b.WriteString("         Typing 'Whole Foods' will suggest 'Groceries'\n")

	return b.String()
}
