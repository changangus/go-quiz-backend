// File: db/migrations/run_migrations.go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Define command-line flags for database connection
	dbUser := flag.String("user", "", "Database user")
	dbPassword := flag.String("password", "", "Database password")
	dbName := flag.String("dbname", "", "Database name")
	dbHost := flag.String("host", "", "Database host")
	dbPort := flag.String("port", "", "Database port")
	migrationsDir := flag.String("migrations", "./db/migrations", "Directory containing migration files")

	// Parse command line flags
	flag.Parse()

	// Use environment variables as defaults if flags are not provided
	if *dbUser == "" {
		*dbUser = getEnv("DB_USER", "postgres")
	}
	if *dbPassword == "" {
		*dbPassword = getEnv("DB_PASSWORD", "password")
	}
	if *dbName == "" {
		*dbName = getEnv("DB_NAME", "quizdb")
	}
	if *dbHost == "" {
		*dbHost = getEnv("DB_HOST", "postgres_quiz_db")
	}
	if *dbPort == "" {
		*dbPort = getEnv("DB_PORT", "5432")
	}

	// PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)

	fmt.Printf("Connecting to PostgreSQL database at %s:%s...\n", *dbHost, *dbPort)

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database")

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}

	// Get list of applied migrations
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		log.Fatal("Failed to query migrations:", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal("Failed to scan migration row:", err)
		}
		appliedMigrations[name] = true
	}

	// Get list of migration files - using os package instead of ioutil
	files, err := os.ReadDir(*migrationsDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	// Filter for .up.sql files and sort them
	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrations = append(migrations, file.Name())
		}
	}

	if len(migrations) == 0 {
		fmt.Println("No migration files found in", *migrationsDir)
		return
	}

	// Run migrations
	for _, migration := range migrations {
		if appliedMigrations[migration] {
			fmt.Println("Migration already applied:", migration)
			continue
		}

		fmt.Println("Applying migration:", migration)

		// Read migration file - using os package instead of ioutil
		content, err := os.ReadFile(filepath.Join(*migrationsDir, migration))
		if err != nil {
			log.Fatal("Failed to read migration file:", err)
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal("Failed to execute migration:", err, "in file:", migration)
		}

		// Record migration
		_, err = db.Exec("INSERT INTO migrations (name) VALUES ($1)", migration)
		if err != nil {
			log.Fatal("Failed to record migration:", err)
		}

		fmt.Println("Successfully applied migration:", migration)
	}

	fmt.Println("All migrations have been applied successfully!")
}
