# ATAD Project - Personal Finance Tracker CLI

A command-line tool for tracking personal income and expenses. Import transactions from bank statements, categorize them automatically, set budgets, and generate insightful reports—all from your terminal.

## Why This Project?

CLI tools force you to think about user experience in a constrained environment. You'll work with file formats, data persistence, and create a practical tool that solves real-world financial tracking needs.

## Current Status: Proof of Concept (PoC)

This is currently a **proof of concept** that demonstrates:
- ✅ Terminal UI application structure
- ✅ Database connectivity
- ✅ Connection testing functionality

## Planned Features

- 📥 Import transactions from bank statements 
- 🏷️ Automatic transaction categorization
- 💰 Budget setting and tracking
- 📊 Financial reports and insights
- 🔍 Transaction search and filtering
- 📈 Spending analytics
- 💾 Data persistence with SQLite

## Tech Stack

- **Language:** Go 1.25+
- **CLI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **Database:** SQLite
- **Database Driver:** [SQLite](https://github.com/mattn/go-sqlite3)

## Prerequisites

- Go 1.25 or higher
- SQLite 

## Setup

### Clone the repository

```bash
git clone https://github.com/PeguB/atad-project.git
cd atad-project
```


### Build and run the application

```bash
go build -o atad ./cmd
```
```bash
./atad 
```