# ATAD CLI Usage Scenario

This document demonstrates a realistic monthly usage scenario using ATAD's CLI subcommands.

## Scenario: December 2025 Personal Finance Management

### Week 1: Setting Up December Budget

**Day 1: Set up monthly budgets for different categories**

```bash
# Set grocery budget for December
./atad budget set Groceries 600 01/12/2025 31/12/2025

# Set entertainment budget
./atad budget set Entertainment 200 01/12/2025 31/12/2025

# Set transportation budget
./atad budget set Transportation 300 01/12/2025 31/12/2025

# Set dining out budget
./atad budget set Restaurants 250 01/12/2025 31/12/2025

# List all budgets to verify
./atad budget list
```

**Expected Output:**
```
✅ Budget set successfully!
   Category: Groceries
   Amount: $600.00
   Period: 01/12/2025 to 31/12/2025

💰 Budgets
──────────────────────────────────────────────────────────────
Category             Amount  Period                  
──────────────────────────────────────────────────────────────
Groceries            $600.00  01/12/2025 - 31/12/2025
Entertainment        $200.00  01/12/2025 - 31/12/2025
Transportation       $300.00  01/12/2025 - 31/12/2025
Restaurants          $250.00  01/12/2025 - 31/12/2025
```

---

### Week 1: Recording Income

**Day 1: Salary arrives**

```bash
# Add monthly salary
./atad add -type income -desc "December Salary" -amount 4500 -category Salary -date 01/12/2025

# Add freelance income
./atad add -type income -desc "Freelance web design project" -amount 800 -category Freelance -date 03/12/2025
```

**Expected Output:**
```
✅ Transaction added successfully!
   Type: Income
   Description: December Salary
   Amount: $4500.00
   Category: Salary
   Date: 01/12/2025
```

---

### Week 1: Daily Expenses

**Day 2: Grocery shopping**

```bash
# Weekly groceries
./atad add -type expense -desc "Whole Foods weekly groceries" -amount 125.50 -date 02/12/2025
```

**Expected Output:**
```
Auto-categorized as: Groceries
✅ Transaction added successfully!
   Type: Expense
   Description: Whole Foods weekly groceries
   Amount: $125.50
   Category: Groceries
   Date: 02/12/2025

💰 Budget: $125.50 / $600.00 (21%)
```

**Day 3: Coffee and lunch**

```bash
# Morning coffee
./atad add -type expense -desc "Starbucks coffee" -amount 6.50 -date 03/12/2025

# Lunch with colleagues
./atad add -type expense -desc "Lunch at Italian bistro" -amount 35.00 -date 03/12/2025
```

**Expected Output:**
```
Auto-categorized as: Restaurants
✅ Transaction added successfully!
   Type: Expense
   Description: Starbucks coffee
   Amount: $6.50
   Category: Restaurants
   Date: 03/12/2025

💰 Budget: $6.50 / $250.00 (3%)
```

**Day 4: Transportation**

```bash
# Gas fill-up
./atad add -type expense -desc "Shell gas station" -amount 55.00 -date 04/12/2025

# Uber to meeting
./atad add -type expense -desc "Uber to downtown meeting" -amount 18.50 -date 04/12/2025
```

---

### Week 2: Mid-Week Check

**Day 8: Check budget status**

```bash
# Check how we're doing on groceries
./atad budget check Groceries

# Check restaurants budget
./atad budget check Restaurants

# Check transportation budget
./atad budget check Transportation
```

**Expected Output:**
```
💰 Budget Status: Groceries
─────────────────────────────────────
Budget:     $600.00
Spent:      $125.50 (21%)
Remaining:  $474.50
Period:     01/12/2025 - 31/12/2025

✅ Within budget
```

**Day 9: More transactions**

```bash
# Second grocery trip
./atad add -type expense -desc "Trader Joe's shopping" -amount 89.30 -date 09/12/2025

# Movie night
./atad add -type expense -desc "Netflix subscription" -amount 15.99 -category Entertainment -date 09/12/2025

# Dinner out
./atad add -type expense -desc "Sushi restaurant dinner" -amount 72.00 -date 09/12/2025
```

---

### Week 3: Searching for Specific Transactions

**Day 15: Need to find all coffee purchases**

```bash
# Search for coffee-related transactions
./atad search "coffee"

# Search for all Uber rides
./atad search "uber"

# Search for groceries
./atad search "grocery"
```

**Expected Output:**
```
🔍 Search Results for 'coffee' (3 found)
─────────────────────────────────────────────────────────────────────────────
Date         Type       Description               Category        Amount
─────────────────────────────────────────────────────────────────────────────
03/12/2025   💸 Expense  Starbucks coffee          Restaurants       6.50
07/12/2025   💸 Expense  Coffee Bean morning...    Restaurants       5.75
12/12/2025   💸 Expense  Local coffee shop         Restaurants       4.25
```

---

### Week 3: Heavy Spending Week

**Day 16-20: Multiple expenses approach budget limits**

```bash
# Big grocery haul
./atad add -type expense -desc "Costco bulk shopping" -amount 185.00 -date 16/12/2025

# Multiple restaurant visits
./atad add -type expense -desc "Birthday dinner celebration" -amount 95.00 -date 17/12/2025
./atad add -type expense -desc "Pizza delivery" -amount 28.50 -date 18/12/2025
./atad add -type expense -desc "Coffee and pastry" -amount 12.00 -date 19/12/2025
./atad add -type expense -desc "Lunch meeting" -amount 42.00 -date 20/12/2025
```

**Expected Output (after birthday dinner):**
```
Auto-categorized as: Restaurants
✅ Transaction added successfully!
   Type: Expense
   Description: Birthday dinner celebration
   Amount: $95.00
   Category: Restaurants
   Date: 17/12/2025

⚠️  Budget warning: $219.00 / $250.00 (88%)
```

---

### Week 4: End of Month Review

**Day 25: List recent transactions**

```bash
# List all transactions from this month
./atad list -limit 50

# List only expenses
./atad list -type expense -limit 30

# List only income
./atad list -type income
```

**Expected Output:**
```
📋 Transactions (30 of 45)
─────────────────────────────────────────────────────────────────────────────
Date         Type       Description               Category        Amount
─────────────────────────────────────────────────────────────────────────────
20/12/2025   💸 Expense  Lunch meeting             Restaurants      42.00
19/12/2025   💸 Expense  Coffee and pastry         Restaurants      12.00
18/12/2025   💸 Expense  Pizza delivery            Restaurants      28.50
17/12/2025   💸 Expense  Birthday dinner celeb...  Restaurants      95.00
16/12/2025   💸 Expense  Costco bulk shopping      Groceries       185.00
...
```

**Day 28: Generate monthly reports**

```bash
# Generate expense report for December
./atad report expense -period month

# Generate income report for December
./atad report income -period month
```

**Expected Output:**
```
📊 Expense Report - December 2025

Total Expense: $1,247.54

Breakdown by Category:
─────────────────────────────────────
  Groceries            $  485.80  (38.9%)
  Restaurants          $  312.25  (25.0%)
  Transportation       $  189.50  (15.2%)
  Entertainment        $   87.99  (7.1%)
  Utilities            $  122.00  (9.8%)
  Healthcare           $   50.00  (4.0%)
```

```
📊 Income Report - December 2025

Total Income: $5,300.00

Breakdown by Category:
─────────────────────────────────────
  Salary               $ 4500.00  (84.9%)
  Freelance            $  800.00  (15.1%)
```

**Day 29: Final budget checks**

```bash
# Check all budgets before month end
./atad budget check Groceries
./atad budget check Restaurants
./atad budget check Transportation
./atad budget check Entertainment
```

**Expected Output (Restaurants - Over Budget!):**
```
💰 Budget Status: Restaurants
─────────────────────────────────────
Budget:     $250.00
Spent:      $312.25 (125%)
Remaining:  $-62.25
Period:     01/12/2025 - 31/12/2025

⚠️  Over budget by $62.25!
```

---

### Week 4: Year-End Analysis

**Day 30: Compare full year performance**

```bash
# Get full year expense report
./atad report expense -period year

# Get full year income report
./atad report income -period year

# Get all-time overview
./atad report expense -period all
./atad report income -period all
```

**Expected Output:**
```
📊 Expense Report - 2025

Total Expense: $48,523.67

Breakdown by Category:
─────────────────────────────────────
  Rent                 $18000.00  (37.1%)
  Groceries            $ 6450.50  (13.3%)
  Transportation       $ 4890.30  (10.1%)
  Restaurants          $ 3245.89  (6.7%)
  Utilities            $ 2456.00  (5.1%)
  Entertainment        $ 1890.45  (3.9%)
  Healthcare           $ 1234.56  (2.5%)
  Shopping             $ 8756.97  (18.0%)
  Education            $ 1600.00  (3.3%)
```

---

## Quick Reference Commands

### Common Daily Usage
```bash
# Add quick expense (auto-categorized)
./atad add -type expense -desc "Lunch" -amount 15

# Check specific budget
./atad budget check Groceries

# Search recent transactions
./atad search "coffee"

# Quick list of today's transactions
./atad list -limit 5
```

### Weekly Review
```bash
# List all transactions
./atad list -limit 50

# Check all budget statuses
./atad budget list

# Generate monthly report
./atad report expense -period month
```

### Monthly Setup
```bash
# Set up new month's budgets
./atad budget set Groceries 600 01/01/2026 31/01/2026
./atad budget set Restaurants 250 01/01/2026 31/01/2026

# Review previous month
./atad report expense -period month
./atad report income -period month
```

---

## Tips for Effective Usage

1. **Use Auto-categorization**: Let ATAD categorize for you when possible
   ```bash
   ./atad add -type expense -desc "Whole Foods" -amount 50
   # Auto-categorized as: Groceries
   ```

2. **Set Budgets at Month Start**: Establish budgets on the 1st of each month
   ```bash
   ./atad budget set Category 500 01/MM/YYYY 31/MM/YYYY
   ```

3. **Regular Budget Checks**: Check budgets weekly to stay on track
   ```bash
   ./atad budget check Groceries
   ```

4. **End-of-Month Reviews**: Generate reports to understand spending patterns
   ```bash
   ./atad report expense -period month
   ```

5. **Search for Patterns**: Find similar transactions to identify habits
   ```bash
   ./atad search "coffee"  # How much am I spending on coffee?
   ```

6. **Use Interactive Mode for Complex Tasks**: Launch full TUI for browsing
   ```bash
   ./atad  # No arguments = Interactive mode
   ```

---

## Combining with Shell Scripts

Create a shell script for recurring tasks:

```bash
#!/bin/bash
# monthly-setup.sh - Set up budgets for new month

MONTH_START="01/01/2026"
MONTH_END="31/01/2026"

echo "Setting up budgets for January 2026..."

./atad budget set Groceries 600 $MONTH_START $MONTH_END
./atad budget set Restaurants 250 $MONTH_START $MONTH_END
./atad budget set Transportation 300 $MONTH_START $MONTH_END
./atad budget set Entertainment 200 $MONTH_START $MONTH_END
./atad budget set Utilities 150 $MONTH_START $MONTH_END

echo "✅ All budgets set!"
./atad budget list
```

---

## Expected Results Summary

After following this scenario, you would have:

- **4 budgets set** for December 2025
- **40+ transactions** recorded across multiple categories
- **Budget warnings** when approaching or exceeding limits
- **Categorized spending** showing patterns
- **Income vs Expense** comparison reports
- **Search capability** to find specific transaction types

This demonstrates how ATAD helps you:
1. ✅ Track daily expenses efficiently
2. ✅ Stay within budget with real-time alerts
3. ✅ Understand spending patterns through reports
4. ✅ Make informed financial decisions
