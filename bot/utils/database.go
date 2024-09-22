// Insert a new user into the users table
package utils

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func DataBaseConnection() {
	dbPath := "../database/UserData.db"

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	fmt.Println("Database file path:", absPath)

	db, err = sql.Open("sqlite", absPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		points INTEGER,
		joinedDate TEXT,
		lastComment TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	fmt.Println("Database and table are set up successfully.")
}
func InsertUser(username string, points int, joinedDate, lastComment string) {
	insertSQL := `INSERT INTO users (username, points, joinedDate, lastComment) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("Error preparing insert statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(username, points, joinedDate, lastComment)
	if err != nil {
		log.Fatalf("Error executing insert statement: %v", err)
	}

	fmt.Println("New user inserted successfully.")
}

// Retrieve all users from the users table
func GetAllUsers() {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("Error retrieving users: %v", err)
	}
	defer rows.Close()

	fmt.Println("Users in the database:")
	for rows.Next() {
		var id int
		var username string
		var points int
		var joinedDate, lastComment string

		err = rows.Scan(&id, &username, &points, &joinedDate, &lastComment)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("ID: %d, Username: %s, Points: %d, JoinedDate: %s, LastComment: %s\n", id, username, points, joinedDate, lastComment)
	}
}

// Update user points based on username
func UpdateUserPoints(username string, newPoints int) {
	updateSQL := `UPDATE users SET points = ? WHERE username = ?`
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		log.Fatalf("Error preparing update statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(newPoints, username)
	if err != nil {
		log.Fatalf("Error executing update statement: %v", err)
	}

	fmt.Printf("User '%s' points updated to %d.\n", username, newPoints)
}
func GetUserPoints(username string) (int, error) {
	var points int
	querySQL := `SELECT points FROM users WHERE username = ?`
	err := db.QueryRow(querySQL, username).Scan(&points)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist
			return 0, nil // or return an appropriate error if preferred
		}
		return 0, fmt.Errorf("error fetching user points: %v", err)
	}
	fmt.Printf("User '%s' points: %d\n", username, points)
	return points, nil
}
func CheckUserExists(username string) (bool, error) {
	var count int
	querySQL := `SELECT COUNT(*) FROM users WHERE username = ?`
	err := db.QueryRow(querySQL, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}
	return count > 0, nil
}
func InsertOrUpdateUser(username, lastComment string) error {
	// Check if the user exists
	exists, err := CheckUserExists(username)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	if exists {
		// If user exists, update their points
		updateSQL := `UPDATE users SET points = points + 1, lastComment = ? WHERE username = ?`
		result, err := db.Exec(updateSQL, lastComment, username)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}
		fmt.Println("User points updated successfully.", result)
	} else {
		// If user doesn't exist, insert a new user with 1 point
		insertSQL := `INSERT INTO users (username, points, joinedDate, lastComment) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertSQL, username, 1, "2023-01-15", lastComment) // Or you can set `joinedDate` to the current date
		if err != nil {
			return fmt.Errorf("error inserting new user: %v", err)
		}
		fmt.Println("New user inserted successfully.")
	}
	return nil
}

// Delete a user by username
func DeleteUser(username string) {
	deleteSQL := `DELETE FROM users WHERE username = ?`
	statement, err := db.Prepare(deleteSQL)
	if err != nil {
		log.Fatalf("Error preparing delete statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(username)
	if err != nil {
		log.Fatalf("Error executing delete statement: %v", err)
	}

	fmt.Printf("User '%s' deleted successfully.\n", username)
}

// Close the database connection
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	fmt.Println("Database connection closed.")
}
