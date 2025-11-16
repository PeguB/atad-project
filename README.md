# ATAD Project - Personal Finance Tracker CLI

A command-line tool for tracking personal income and expenses. Import transactions from bank statements, categorize them automatically, set budgets, and generate insightful reports—all from your terminal.

## Why This Project?

CLI tools force you to think about user experience in a constrained environment. You'll work with file formats, data persistence, and create a practical tool that solves real-world financial tracking needs.

## Current Status: Proof of Concept (PoC)

This is currently a **proof of concept** that demonstrates:
- ✅ CLI application structure using Cobra
- ✅ Database connectivity (PostgreSQL)
- ✅ Connection testing functionality

## Planned Features

- 📥 Import transactions from bank statements (CSV, JSON)
- 🏷️ Automatic transaction categorization
- 💰 Budget setting and tracking
- 📊 Financial reports and insights
- 🔍 Transaction search and filtering
- 📈 Spending analytics
- 💾 Data persistence with PostgreSQL

## Tech Stack

- **Language:** Go 1.25+
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
- **Database:** PostgreSQL
- **Database Driver:** [lib/pq](https://github.com/lib/pq)

## Prerequisites

- Go 1.25 or higher
- Docker & Docker Compose (for local database)
- PostgreSQL (if not using Docker)

## Setup

### Clone the repository

```bash
git clone https://github.com/PeguB/atad-project.git
cd atad-project
```



## Usage

### Build the application

```bash
go build -o atad ./cmd
```

### Run commands

**Test database connection:**
```bash
./atad test-db
```

Or run directly with Go:
```bash
go run ./cmd test-db
```

**Show available commands:**
```bash
./atad --help
```

## Project Structure

```
.
├── cmd/                    # CLI commands
│   ├── main.go            # Application entry point
│   ├── root.go            # Root command and DB connection
│   └── test_db.go         # Database connection test command
├── internal/
│   ├── database/          # Database connection logic
│   ├── models/            # Data models (to be implemented)
│   └── repository/        # Database repositories (to be implemented)
├── docs/                  # Documentation
├── tests/                 # Tests
├── docker-compose.yml     # PostgreSQL setup
├── .env.example           # Environment variables template
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Author

Bogdan Pegulescu (@PeguB)
