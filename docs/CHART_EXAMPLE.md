# Terminal Charts Feature

ATAD now includes visual ASCII bar charts in report output!

## Example Output

When you run a report command like:

```bash
./atad report expense -period month
```

You'll now see a visual bar chart showing category breakdown:

```
📊 Expense Report - December 2025

Total Expense: $450.00

Category Breakdown Chart:
─────────────────────────────────────────────────────────────────────────────
Groceries            │████████████████████████████████████████████████ $150.00 (33.3%)
Restaurants          │████████████████████████████████ $120.00 (26.7%)
Transportation       │████████████████████ $80.00 (17.8%)
Entertainment        │████████████ $50.00 (11.1%)
Coffee               │████████ $30.00 (6.7%)
Shopping             │████ $20.00 (4.4%)
─────────────────────────────────────────────────────────────────────────────

Breakdown by Category:
─────────────────────────────────────
  Groceries            $  150.00  (33.3%)
  Restaurants          $  120.00  (26.7%)
  Transportation       $   80.00  (17.8%)
  Entertainment        $   50.00  (11.1%)
  Coffee               $   30.00  (6.7%)
  Shopping             $   20.00  (4.4%)
```

## Features

- **Visual Bar Chart**: Each category gets a proportional bar made of `█` characters
- **Sorted by Amount**: Categories are automatically sorted from highest to lowest spending
- **Percentage Display**: Shows both dollar amount and percentage of total
- **Auto-scaling**: Bars scale relative to the largest category
- **Clean Layout**: Aligned columns for easy reading

## How It Works

The chart functionality:
1. Collects all transactions for the specified period
2. Groups them by category
3. Sorts categories by spending amount (highest first)
4. Calculates proportional bar widths (max 50 characters)
5. Renders each bar with amount and percentage

This makes it easy to see at a glance where your money is going! 📊
