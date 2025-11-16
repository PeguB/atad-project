package main

import (
	"log"

	"github.com/spf13/cobra"
)

var testDBCmd = &cobra.Command{
	Use:   "test-db",
	Short: "Test the database connection",
	Long:  `Attempts to connect to the database and verifies the connection is working`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Testing database connection...")

		db, err := getDB()
		if err != nil {
			log.Fatalf("❌ Failed to connect to database: %v", err)
		}

		// Test the connection with a ping
		if err := db.DB.Ping(); err != nil {
			log.Fatalf("❌ Database connection failed: %v", err)
		}

		log.Println("✅ Database connection successful!")
	},
}

func init() {
	rootCmd.AddCommand(testDBCmd)
}
