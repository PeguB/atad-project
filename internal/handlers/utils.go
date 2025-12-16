package handlers

import (
	"fmt"
	"sort"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

// TruncateString truncates a string to maxLen, adding "..." if needed
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// CategoryColor holds category name and its assigned color
type CategoryColor struct {
	Category string
	Color    lipgloss.Color
	Style    lipgloss.Style
}

// categoryColors is a package-level variable to store color assignments
var categoryColors = make(map[string]CategoryColor)

// DrawCategoryBarChart renders a bar chart for category spending using ntcharts
func DrawCategoryBarChart(byCategory map[string]float64, total float64) {
	if len(byCategory) == 0 {
		return
	}

	// Sort categories by amount
	type categoryAmount struct {
		category string
		amount   float64
	}
	var sorted []categoryAmount
	for cat, amt := range byCategory {
		sorted = append(sorted, categoryAmount{cat, amt})
	}

	// Sort descending by amount
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].amount > sorted[j].amount
	})

	// Define color palette for bars
	colorPalette := []lipgloss.Color{
		lipgloss.Color("10"),  // bright green
		lipgloss.Color("14"),  // bright cyan
		lipgloss.Color("11"),  // bright yellow
		lipgloss.Color("13"),  // bright magenta
		lipgloss.Color("9"),   // bright red
		lipgloss.Color("12"),  // bright blue
		lipgloss.Color("208"), // orange
		lipgloss.Color("165"), // purple
	}

	// Prepare data for bar chart
	barData := make([]barchart.BarData, 0, len(sorted))

	for i, item := range sorted {
		// Assign color to category
		color := colorPalette[i%len(colorPalette)]
		barStyle := lipgloss.NewStyle().
			Foreground(color).
			Background(color)

		// Store color assignment for later use in breakdown
		categoryColors[item.category] = CategoryColor{
			Category: item.category,
			Color:    color,
			Style:    barStyle,
		}

		// Truncate long category names
		displayCategory := item.category
		if len(displayCategory) > 15 {
			displayCategory = displayCategory[:12] + "..."
		}

		barData = append(barData, barchart.BarData{
			Label: displayCategory,
			Values: []barchart.BarValue{
				{
					Value: item.amount,
					Style: barStyle,
				},
			},
		})
	}

	// Create bar chart with horizontal bars
	chartHeight := len(sorted)*2 + 4
	if chartHeight > 30 {
		chartHeight = 30
	}

	// Find max value for proper scaling
	maxValue := 0.0
	for _, item := range sorted {
		if item.amount > maxValue {
			maxValue = item.amount
		}
	}

	// Create bar chart with horizontal bars and labels
	axisStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))   // white
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // cyan

	bc := barchart.New(80, chartHeight,
		barchart.WithHorizontalBars(),
		barchart.WithDataSet(barData),
		barchart.WithMaxValue(maxValue*1.1), // Add 10% padding
		barchart.WithStyles(axisStyle, labelStyle),
	)

	// Draw the chart before rendering
	bc.Draw()

	// Render the chart
	fmt.Println("\nCategory Breakdown Chart:")
	fmt.Println(bc.View())
}

// GetCategoryColor returns the color style for a category
func GetCategoryColor(category string) lipgloss.Style {
	if cc, ok := categoryColors[category]; ok {
		return cc.Style
	}
	// Default style if color not found
	return lipgloss.NewStyle()
}
