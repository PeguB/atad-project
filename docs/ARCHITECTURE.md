# Command Handler Architecture

This document describes the refactored CLI command handler architecture.

## Overview

The CLI command handlers have been refactored into a clean, maintainable structure using interfaces and are organized within the `internal/handlers` package for better code organization and encapsulation.

## File Structure

```
cmd/
└── main.go                 # TUI and main entry point

internal/
├── handlers/
│   ├── cli_handlers.go    # CLI command handlers
│   └── utils.go           # Helper functions (TruncateString, DrawCategoryBarChart)
├── database/
├── models/
├── repository/
├── service/
└── tui/
```

## Architecture

### CommandHandler Interface

All command handlers implement the `CommandHandler` interface:

```go
type CommandHandler interface {
    Handle()
}
```

### CLIHandler Struct

The `CLIHandler` is the main struct that manages database connections and repositories:

```go
type CLIHandler struct {
    db              *database.Database
    txRepo          *repository.TransactionRepository
    budgetRepo      *repository.BudgetRepository
    categoryService *service.CategoryService
}
```

**Methods:**
- `NewCLIHandler()` - Creates a new CLI handler instance
- `initDatabase()` - Initializes database connection and repositories
- `Close()` - Closes the database connection

### Command Structs

Each command is implemented as a separate struct that holds a reference to the CLIHandler:

1. **AddCommand** - Handles `atad add` command
2. **ListCommand** - Handles `atad list` command
3. **ReportCommand** - Handles `atad report` command
4. **BudgetCommand** - Handles `atad budget` command
   - `handleList()` - Lists all budgets
   - `handleSet()` - Sets a new budget
   - `handleCheck()` - Checks budget status
5. **SearchCommand** - Handles `atad search` command
6. **ImportCommand** - Handles `atad import` command (placeholder)

## Benefits of This Architecture

1. **Separation of Concerns**: TUI code in `cmd/main.go`, CLI commands in `internal/handlers/`
2. **Encapsulation**: Handlers are in the internal package, preventing external imports
3. **Testability**: Each command can be tested independently
4. **Maintainability**: Easy to add new commands or modify existing ones
5. **Single Responsibility**: Each struct has one clear purpose
6. **Resource Management**: Database connections are properly managed through the CLIHandler
7. **Interface-based Design**: Easy to mock or extend functionality
8. **Package Organization**: Following Go best practices with internal packages

## Usage Flow

1. User runs a command: `./atad report expense -period month`
2. `main()` creates a `handlers.CLIHandler` instance
3. Command is routed to appropriate handler (e.g., `handlers.ReportCommand`)
4. Handler calls `initDatabase()` to set up connections
5. Handler executes its `Handle()` method
6. Database connection is closed via `defer Close()`

## Adding a New Command

To add a new command:

1. Create a new struct in `internal/handlers/cli_handlers.go`:
```go
type NewCommand struct {
    Handler *CLIHandler
}
```

2. Implement the `Handle()` method:
```go
func (c *NewCommand) Handle() {
    // Command logic here
}
```

3. Add routing in `cmd/main.go`:
```go
case "newcommand":
    cmd = &handlers.NewCommand{Handler: handler}
```

## Example: How ReportCommand Works

```go
// User runs: ./atad report expense -period month

// 1. main() creates handler
handler := handlers.NewCLIHandler()

// 2. Creates command instance
cmd := &handlers.ReportCommand{Handler: handler}

// 3. Executes command
cmd.Handle()
    ├── Initializes database
    ├── Parses command flags
    ├── Retrieves transactions
    ├── Calculates totals by category
    ├── Draws chart (from handlers.DrawCategoryBarChart)
    └── Displays results
```

## Utility Functions

Located in `internal/handlers/utils.go`:

- **TruncateString(s string, maxLen int)** - Truncates strings with ellipsis
- **DrawCategoryBarChart(byCategory map[string]float64, total float64)** - Renders ASCII bar charts

These utilities are exported (capitalized) so they can be used across the handlers package while still being internal to the application.

## Package Visibility

- **Public (Exported)**: `CommandHandler`, `CLIHandler`, `NewCLIHandler`, command structs, utility functions
- **Private (Unexported)**: `initDatabase()`, `Close()`, and private helper methods like `handleList()`, `handleSet()`, `handleCheck()`

This visibility design ensures that only the necessary interfaces and constructors are exposed while keeping implementation details private.
