package database

import (
	"classplanner/internal/repository"
	"classplanner/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Import pq for PostgreSQL driver
)

type DatabaseInstance struct {
	DB         *sql.DB
	Repository *repository.Repository
}

// database instance that contains the conection & the repository
var DBInstance *DatabaseInstance

func Connect() *DatabaseInstance {
	// load env variables
	utils.LoadEnv()

	// Define the connection string with PostgreSQL credentials
	user := os.Getenv("DATABASE_USER")
	passwd := os.Getenv("DATABASE_PASSWORD")
	database := os.Getenv("DATABASE_DATABASE")
	ssl := os.Getenv("DATABASE_SSL")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, passwd, database, ssl)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	DBInstance = &DatabaseInstance{DB: db}
	DBInstance.Repository = repository.New(db)

	log.Println("✅ Database connected successfully")

	return DBInstance
}

// Ready checks DB connection health
func (r *DatabaseInstance) Ready() bool {
	if r == nil || r.DB == nil {
		return false
	}
	if err := r.DB.Ping(); err != nil {
		log.Printf("⚠️ Database not ready: %v", err)
		return false
	}
	return true
}
