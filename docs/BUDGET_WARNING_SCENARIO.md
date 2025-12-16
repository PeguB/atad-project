# Budget Warning Scenario

This scenario demonstrates ATAD's budget warning system when spending approaches or exceeds budget limits.

## Scenario: Coffee Budget Gone Wrong ☕

### Step 1: Set a Small Coffee Budget

```bash
# Set a modest $50 coffee budget for December
./atad budget set Coffee 50 01/12/2025 31/12/2025
```

**Output:**
```
✅ Budget set successfully!
   Category: Coffee
   Amount: $50.00
   Period: 01/12/2025 to 31/12/2025
```

---

### Step 2: First Few Purchases (Under Budget)

```bash
# Day 1: Morning coffee
./atad add -type expense -desc "Starbucks latte" -amount 6.50 -category Coffee -date 01/12/2025
```

**Output:**
```
✅ Transaction added successfully!
   Type: Expense
   Description: Starbucks latte
   Amount: $6.50
   Category: Coffee
   Date: 01/12/2025

💰 Budget: $6.50 / $50.00 (13%)
```

```bash
# Day 3: Another coffee
./atad add -type expense -desc "Local cafe cappuccino" -amount 5.75 -category Coffee -date 03/12/2025
```

**Output:**
```
✅ Transaction added successfully!
   Type: Expense
   Description: Local cafe cappuccino
   Amount: $5.75
   Category: Coffee
   Date: 03/12/2025

💰 Budget: $12.25 / $50.00 (25%)
```

---

### Step 3: Approaching Warning Threshold (80%)

```bash
# Day 5: Coffee and pastry
./atad add -type expense -desc "Coffee Bean with pastry" -amount 12.00 -category Coffee -date 05/12/2025

# Day 7: Premium coffee
./atad add -type expense -desc "Specialty pour-over" -amount 8.50 -category Coffee -date 07/12/2025

# Day 9: Meeting at cafe
./atad add -type expense -desc "Cafe meeting - 2 lattes" -amount 13.00 -category Coffee -date 09/12/2025
```

**Output (after Day 9):**
```
✅ Transaction added successfully!
   Type: Expense
   Description: Cafe meeting - 2 lattes
   Amount: $13.00
   Category: Coffee
   Date: 09/12/2025

⚠️  Budget warning: $45.75 / $50.00 (92%)
```

---

### Step 4: Exceeding Budget! 🚨

```bash
# Day 11: Forgot about the budget...
./atad add -type expense -desc "Starbucks grande mocha" -amount 7.25 -category Coffee -date 11/12/2025
```

**Output:**
```
✅ Transaction added successfully!
   Type: Expense
   Description: Starbucks grande mocha
   Amount: $7.25
   Category: Coffee
   Date: 11/12/2025

⚠️  OVER BUDGET: $53.00 / $50.00 (106%)
You've exceeded your budget by $3.00
```

---

### Step 5: Keep Going (Making It Worse)

```bash
# Day 13: Another purchase
./atad add -type expense -desc "Coffee to-go" -amount 5.50 -category Coffee -date 13/12/2025
```

**Output:**
```
✅ Transaction added successfully!
   Type: Expense
   Description: Coffee to-go
   Amount: $5.50
   Category: Coffee
   Date: 13/12/2025

⚠️  OVER BUDGET: $58.50 / $50.00 (117%)
You've exceeded your budget by $8.50
```

---

### Step 6: Check Budget Status

```bash
# Check how bad it really is
./atad budget check Coffee
```

**Output:**
```
💰 Budget Status: Coffee
─────────────────────────────────────
Budget:     $50.00
Spent:      $58.50 (117%)
Remaining:  $-8.50
Period:     01/12/2025 - 31/12/2025

⚠️  Over budget by $8.50!
```

---

### Step 7: View All Coffee Transactions

```bash
# See where all the money went
./atad search "coffee"
```

**Output:**
```
🔍 Search Results for 'coffee' (7 found)
─────────────────────────────────────────────────────────────────────────────
Date         Type       Description               Category        Amount
─────────────────────────────────────────────────────────────────────────────
13/12/2025   💸 Expense  Coffee to-go              Coffee            5.50
11/12/2025   💸 Expense  Starbucks grande mocha    Coffee            7.25
09/12/2025   💸 Expense  Cafe meeting - 2 lattes   Coffee           13.00
07/12/2025   💸 Expense  Specialty pour-over       Coffee            8.50
05/12/2025   💸 Expense  Coffee Bean with pastry   Coffee           12.00
03/12/2025   💸 Expense  Local cafe cappuccino     Coffee            5.75
01/12/2025   💸 Expense  Starbucks latte           Coffee            6.50
─────────────────────────────────────────────────────────────────────────────
Total: $58.50
```

---

## Summary of Warning Levels

| Spending | Percentage | Warning Level |
|----------|-----------|---------------|
| $12.25   | 25%       | ✅ Normal - Within budget |
| $45.75   | 92%       | ⚠️  Warning - Approaching limit |
| $53.00   | 106%      | 🚨 Over Budget - Exceeded by $3.00 |
| $58.50   | 117%      | 🚨 Over Budget - Exceeded by $8.50 |

---

## Key Takeaways

1. **80% Threshold**: Warning appears when you reach 80% of budget
2. **100% Exceeded**: Red alert when you go over budget
3. **Real-time Feedback**: Every transaction shows current budget status
4. **Easy Tracking**: `budget check` command shows exact overage amount
5. **Pattern Analysis**: Search helps identify spending habits

This scenario shows how ATAD helps you stay aware of your spending in real-time! 💰
