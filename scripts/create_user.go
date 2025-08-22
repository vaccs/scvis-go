package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	// Load environment variables - try scvisual.conf first, then fallback to .env
	err := godotenv.Load("/etc/scvisual/scvisual.conf")
	if err != nil {
		// Fallback to .env if scvisual.conf doesn't exist
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading configuration file (.env or scvisual.conf)")
		}
	}

	// Get database connection details from environment
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		dbHost = "localhost"
	}

	if dbUser == "" {
		log.Fatal("DB_USER not found in configuration")
	}

	if dbPasswd == "" {
		log.Fatal("DB_PASSWD not found in configuration")
	}

	if dbName == "" {
		log.Fatal("DB_NAME not found in configuration")
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPasswd, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to database successfully!")
	fmt.Println("Creating new user for scvis_local database")
	fmt.Println("=========================================")

	// Get user input
	reader := bufio.NewReader(os.Stdin)

	// Get name
	fmt.Print("Enter user name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading name:", err)
	}
	name = strings.TrimSpace(name)

	// Get email
	fmt.Print("Enter email address: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading email:", err)
	}
	email = strings.TrimSpace(email)

	// Get password (hidden input)
	fmt.Print("Enter password (min 8 characters): ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("Error reading password:", err)
	}
	password := string(passwordBytes)
	fmt.Println() // New line after hidden input

	// Confirm password (hidden input)
	fmt.Print("Confirm password: ")
	confirmPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("Error reading password confirmation:", err)
	}
	confirmPassword := string(confirmPasswordBytes)
	fmt.Println() // New line after hidden input

	// Validate inputs
	if name == "" {
		log.Fatal("Name cannot be empty")
	}

	if email == "" {
		log.Fatal("Email cannot be empty")
	}

	if len(password) < 8 {
		log.Fatal("Password must be at least 8 characters long")
	}

	if password != confirmPassword {
		log.Fatal("Passwords do not match")
	}

	// Get administrator status
	fmt.Print("Is this user an administrator? (y/N): ")
	adminInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading administrator input:", err)
	}
	adminInput = strings.TrimSpace(strings.ToLower(adminInput))
	isAdmin := adminInput == "y" || adminInput == "yes"

	// Hash password using bcrypt with cost 12 (same as the application)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	// Insert user into database
	stmt := `INSERT INTO users (name, email, hashed_password, administrator)
             VALUES(?, ?, ?, ?)`

	_, err = db.Exec(stmt, name, email, string(hashedPassword), isAdmin)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "email") {
			log.Fatal("Error: A user with this email address already exists")
		}
		log.Fatal("Error inserting user:", err)
	}

	fmt.Printf("\nUser created successfully!\n")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Administrator: %t\n", isAdmin)
	fmt.Println("\nThe user can now log in to the application.")
}
