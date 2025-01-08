package main

import (
	"database/sql"
	"fmt"
	"log"
)

type UserRepository interface {
	GetUserByID(id int) (string, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// GetUserByID retrieves a user by ID from the PostgreSQL database
func (r *PostgresUserRepository) GetUserByID(id int) (string, error) {
	var username string
	query := "SELECT username FROM users WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

type UserService struct {
	repo UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUser retrieves a user by ID using the UserRepository
func (s *UserService) GetUser(id int) (string, error) {
	return s.repo.GetUserByID(id)
}

// OpenPostgresDBConnection open connection to database
func OpenPostgresDBConnection() (*sql.DB, error) {

	// Database connection string (PostgreSQL example)
	connStr := "user=youruser password=yourpassword dbname=yourdb sslmode=disable"

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	return db, nil
}

func main() {
	// Open a database connection
	db, err := OpenPostgresDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Create a PostgresUserRepository instance
	userRepo := NewPostgresUserRepository(db)

	// Inject the userRepo into UserService
	userService := NewUserService(userRepo)

	// Fetch a user by ID
	user, err := userService.GetUser(1)
	if err != nil {
		log.Fatalf("Error retrieving user: %v", err)
	}

	fmt.Printf("Retrieved user: %s\n", user)
}
