package utils

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

type User struct {
	ID          int
	UserName    string
	UserId      string
	Points      int
	JoinedDate  string
	LastComment string
	LastSeen    string
	ProfilePic  string
	CarUrl      string
}
type Product struct {
	Code  string
	Price uint
}

const (
	UserName    = "USER_NAME"
	UserId      = "USER_ID"
	Points      = "POINTS"
	JoinedDate  = "JOINED_DATE"
	LastComment = "LAST_COMMENT"
	LastSeen    = "LAST_SEEN"
	ProfilePic  = "PROFILE_PIC"
)

func DataBaseConnection() {
	dbPath := "./database/UserData.db"

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
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		USER_NAME TEXT NOT NULL,
		USER_ID	TEXT NOT NULL,
		POINTS INTEGER,
		JOINED_DATE TEXT,
		LAST_COMMENT TEXT,
		LAST_SEEN TEXT,
		PROFILE_PIC TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	fmt.Println("Database and table are set up successfully.")
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

// Update user points based on userid
func UpdateUserPoints(authorId string, newPoints int) {
	updateSQL := `UPDATE users SET points = ? WHERE userId = ?`
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		log.Fatalf("Error preparing update statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(newPoints, authorId)
	if err != nil {
		log.Fatalf("Error executing update statement: %v", err)
	}

	fmt.Printf("User '%s' points updated to %d.\n", authorId, newPoints)
}

func CheckUserExists(msg ChatMessage) (bool, error) {
	var count int

	querySQL := fmt.Sprintf(`SELECT COUNT(*) FROM users WHERE %s = ?`, UserId)
	err := db.QueryRow(querySQL, msg.AuthorId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}
	return count > 0, nil
}
func InsertOrUpdateUser(msg ChatMessage) error {
	// Check if the user exists
	fmt.Println("Checking if exits ")

	exists, err := CheckUserExists(msg)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}
	fmt.Println("Response of ", exists)
	if exists {
		// If user exists, update their points
		updateSQL := fmt.Sprintf(`
		UPDATE users 
		SET %s = %s + 1, %s = ?, %s = ? 
		WHERE %s = ?`,
			Points, Points, LastComment, LastSeen, UserId)
		result, err := db.Exec(updateSQL, msg.MessageContent, msg.CommentTime, msg.AuthorId)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}
		fmt.Println("User points updated successfully.", result)
	} else {
		// If user doesn't exist, insert a new user with 1 point
		fmt.Println("Trying to insert a new user")
		insertSQL := fmt.Sprintf(`
		INSERT INTO users (%s, %s, %s, %s, %s, %s,%s) 
		VALUES (?, ?, ?, ?, ?, ?,?)`,
			UserName, UserId, Points, JoinedDate, LastComment, LastSeen, ProfilePic)
		_, err := db.Exec(insertSQL, msg.AuthorName, msg.AuthorId, 1, msg.CommentTime, msg.MessageContent, msg.CommentTime, msg.AuthorPhotoURL) // Or you can set `joinedDate` to the current date
		if err != nil {
			fmt.Println("Error inserting new user: ", err)
			return fmt.Errorf("error inserting new user: %v", err)
		}
		fmt.Println("New user inserted successfully.")
	}
	FetchUser(msg)
	return nil
}

func FetchUser(msg ChatMessage) (User, error) {
	var user User
	querySQL := fmt.Sprintf(`SELECT * FROM users WHERE %s = ?`, UserId)
	err := db.QueryRow(querySQL, msg.AuthorId).Scan(&user.ID, &user.UserName, &user.UserId, &user.Points, &user.JoinedDate, &user.LastComment, &user.LastSeen, &user.ProfilePic)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist
			return user, err
		}
		fmt.Println("Error fetching user:", err)
		return user, err
	}
	return user, nil
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

func Transform(msg ChatMessage) User {

	user := User{
		UserName:    msg.AuthorName,
		UserId:      msg.AuthorId,
		LastComment: msg.MessageContent,
		LastSeen:    msg.CommentTime,
		Points:      0,
		JoinedDate:  msg.CommentTime,
		ProfilePic:  msg.AuthorPhotoURL,
	}

	return user
}
