package models

type CategoryRule struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Pattern     string `json:"pattern"`     // Regex pattern
	Description string `json:"description"` // What this rule matches
	Priority    int    `json:"priority"`    // Higher priority rules are checked first
}

// DefaultRules returns a set of common categorization rules
func DefaultCategoryRules() []CategoryRule {
	return []CategoryRule{
		{Category: "Groceries", Pattern: `(?i)(grocery|supermarket|whole foods|trader joe|safeway|walmart|kroger|costco|food market)`, Description: "Grocery stores", Priority: 10},
		{Category: "Restaurants", Pattern: `(?i)(restaurant|cafe|coffee|starbucks|mcdonald|burger|pizza|diner|bistro|bar & grill)`, Description: "Dining out", Priority: 10},
		{Category: "Transportation", Pattern: `(?i)(uber|lyft|taxi|gas station|shell|chevron|bp|exxon|mobil|parking|metro|transit)`, Description: "Transportation and fuel", Priority: 10},
		{Category: "Utilities", Pattern: `(?i)(electric|water|gas company|utility|internet|phone|wireless|at&t|verizon|comcast)`, Description: "Utility bills", Priority: 10},
		{Category: "Entertainment", Pattern: `(?i)(netflix|spotify|hulu|disney|hbo|cinema|movie|theater|concert|game|steam)`, Description: "Entertainment services", Priority: 10},
		{Category: "Shopping", Pattern: `(?i)(amazon|ebay|target|best buy|apple store|mall|clothing|fashion|retail)`, Description: "General shopping", Priority: 5},
		{Category: "Healthcare", Pattern: `(?i)(pharmacy|cvs|walgreens|hospital|clinic|doctor|dental|medical|health)`, Description: "Medical expenses", Priority: 10},
		{Category: "Fitness", Pattern: `(?i)(gym|fitness|yoga|sports|athletic)`, Description: "Fitness and sports", Priority: 10},
		{Category: "Salary", Pattern: `(?i)(salary|payroll|wages|income|direct deposit)`, Description: "Employment income", Priority: 10},
		{Category: "Investment", Pattern: `(?i)(dividend|interest|capital gain|stock|investment)`, Description: "Investment returns", Priority: 10},
		{Category: "Rent", Pattern: `(?i)(rent|lease|housing)`, Description: "Housing rent", Priority: 10},
		{Category: "Insurance", Pattern: `(?i)(insurance|premium)`, Description: "Insurance payments", Priority: 10},
		{Category: "Education", Pattern: `(?i)(tuition|school|university|course|textbook|education)`, Description: "Educational expenses", Priority: 10},
		{Category: "Travel", Pattern: `(?i)(hotel|airline|booking|airbnb|flight|vacation)`, Description: "Travel expenses", Priority: 10},
	}
}
