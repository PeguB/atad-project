package main

import (
	"log"
	"os"

	"github.com/PeguB/atad-project/internal/database"
	"github.com/spf13/cobra"
)

var db *database.Database

var rootCmd = &cobra.Command{
	Use:   "atad",
	Short: "ATAD Project CLI",
	Long:  `A CLI application for the ATAD project with database connectivity`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// getDB returns a database connection, initializing it if needed
func getDB() (*database.Database, error) {
	if db == nil {
		var err error
		db, err = database.NewDatabase()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
