package service

import (
	"regexp"
	"sort"
	"strings"

	"github.com/PeguB/atad-project/internal/models"
)

type CategoryService struct {
	rules []models.CategoryRule
}

func NewCategoryService() *CategoryService {
	rules := models.DefaultCategoryRules()
	// Sort by priority (highest first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	return &CategoryService{
		rules: rules,
	}
}

// CategorizeTransaction attempts to categorize a transaction based on its description
func (s *CategoryService) CategorizeTransaction(description string) string {
	description = strings.TrimSpace(description)

	// Try to match against each rule
	for _, rule := range s.rules {
		matched, err := regexp.MatchString(rule.Pattern, description)
		if err != nil {
			continue // Skip invalid regex patterns
		}

		if matched {
			return rule.Category
		}
	}

	// Default category if no match found
	return "Uncategorized"
}

// AddCustomRule adds a new categorization rule
func (s *CategoryService) AddCustomRule(category, pattern, description string, priority int) error {
	// Validate the regex pattern
	_, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	rule := models.CategoryRule{
		Category:    category,
		Pattern:     pattern,
		Description: description,
		Priority:    priority,
	}

	s.rules = append(s.rules, rule)

	// Re-sort by priority
	sort.Slice(s.rules, func(i, j int) bool {
		return s.rules[i].Priority > s.rules[j].Priority
	})

	return nil
}

// GetRules returns all categorization rules
func (s *CategoryService) GetRules() []models.CategoryRule {
	return s.rules
}

// TestRule tests a pattern against a description
func (s *CategoryService) TestRule(pattern, description string) (bool, error) {
	return regexp.MatchString(pattern, description)
}
