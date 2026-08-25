package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aaripurna/potash/config"
	"github.com/aaripurna/potash/database"
	"github.com/aaripurna/potash/migrations"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database schema migrations",
	Long: `Run goose migrations against DATABASE_URL.

	migrate up        apply every pending migration
	migrate down      roll the most recent migration back
	migrate status    show which migrations have run
	migrate create    scaffold a new migration file
	`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply every pending migration",
	Run: func(cmd *cobra.Command, args []string) {
		withGoose(func(db *gorm.DB) error { return runGoose(db, goose.Up) })
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll the most recent migration back",
	Run: func(cmd *cobra.Command, args []string) {
		withGoose(func(db *gorm.DB) error { return runGoose(db, goose.Down) })
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show which migrations have run",
	Run: func(cmd *cobra.Command, args []string) {
		withGoose(func(db *gorm.DB) error { return runGoose(db, goose.Status) })
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create NAME",
	Short: "Scaffold a new migration file in ./migrations",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Written to the working tree, so this one is for development only.
		if err := goose.Create(nil, "./migrations", args[0], "sql"); err != nil {
			log.Fatalf("unable to create migration: %v", err)
		}
	},
}

// withGoose opens a connection, hands it to fn, and always closes the pool.
func withGoose(fn func(db *gorm.DB) error) {
	config.InitEnv()

	db, err := database.Open(config.DatabaseURL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("unable to reach the connection pool: %v", err)
	}

	defer sqlDB.Close()

	if err := fn(db); err != nil {
		log.Fatalf("%v", err)
	}
}

func runGoose(db *gorm.DB, action func(*sql.DB, string, ...goose.OptionsFunc) error) error {
	sqlDB, err := db.DB()

	if err != nil {
		return fmt.Errorf("unable to reach the connection pool: %w", err)
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("unable to select the goose dialect: %w", err)
	}

	return action(sqlDB, ".")
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, migrateStatusCmd, migrateCreateCmd)
	rootCmd.AddCommand(migrateCmd)
}
