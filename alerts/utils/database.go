package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

type User struct {
	UserName       string
	UserId         string
	Points         int
	JoinedDate     string
	FirstVideoLink string
	LastVideoLink  string
	LastComment    string
	LastSeen       string
	ProfilePic     string
}

// Convert Go struct (User) to JSON ([]byte)
func ObjectToJSON(user *User) ([]byte, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling user data: %v", err)
	}
	return b, nil
}

// Convert JSON ([]byte) to Go struct (User)
func JSONToObject(data []byte) (User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return user, fmt.Errorf("Error unmarshalling JSON data: %v", err)
	}
	return user, nil
}

func DataBaseConnection() {
	var err error
	// Open an SQLite database
	db, err = sql.Open("sqlite", "./userDB.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create the users table with the new schema
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		UserId TEXT PRIMARY KEY,
		UserJson TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func GetUserById(userId string) (User, error) {
	var user User
	var userJson []byte

	// Query to fetch the UserJson for the specified UserId
	query := `SELECT UserJson FROM users WHERE UserId = ?`
	row := db.QueryRow(query, userId)

	// Scan the JSON data from the database
	err := row.Scan(&userJson)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error retrieving user data: %v", err)
	}

	// Unmarshal the JSON data into the User struct
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return user, fmt.Errorf("error unmarshalling user data: %v", err)
	}

	return user, nil
}

func InsertUser(msg ChatMessage) (User, error) {
	insertUser := User{
		UserName:       msg.AuthorName,
		UserId:         msg.AuthorId,
		LastComment:    msg.MessageContent,
		LastSeen:       msg.CommentTime,
		Points:         1,
		JoinedDate:     msg.CommentTime,
		ProfilePic:     msg.AuthorPhotoURL,
		FirstVideoLink: msg.VideoID,
		LastVideoLink:  msg.VideoID,
	}

	// Convert User struct to JSON
	userJson, err := json.Marshal(insertUser)
	if err != nil {
		return insertUser, fmt.Errorf("error marshalling user data to JSON: %v", err)
	}

	// Insert the user and their JSON data
	query := `INSERT OR REPLACE INTO users (UserId, UserJson) VALUES (?, ?)`
	_, err = db.Exec(query, insertUser.UserId, userJson)
	if err != nil {
		return insertUser, fmt.Errorf("error inserting or updating user: %v", err)
	}

	fmt.Println("User inserted successfully.")
	return insertUser, nil
}
func UpdateUser(msg ChatMessage, user User) (User, error) {
	// Update the user fields
	updatedUser := User{
		UserName:       msg.AuthorName,
		UserId:         msg.AuthorId,
		LastComment:    msg.MessageContent,
		LastSeen:       msg.CommentTime,
		Points:         user.Points + 1,
		JoinedDate:     user.JoinedDate,
		ProfilePic:     msg.AuthorPhotoURL,
		FirstVideoLink: user.FirstVideoLink,
		LastVideoLink:  msg.VideoID,
	}

	// Convert updatedUser to JSON
	userJson, err := json.Marshal(updatedUser)
	if err != nil {
		return updatedUser, fmt.Errorf("error marshalling updated user data to JSON: %v", err)
	}

	// Update the UserJson column in the database
	query := `UPDATE users SET UserJson = ? WHERE UserId = ?`
	_, err = db.Exec(query, userJson, updatedUser.UserId)
	if err != nil {
		return updatedUser, fmt.Errorf("error updating user: %v", err)
	}

	fmt.Println("User updated successfully.")
	return updatedUser, nil
}
func CheckUserExists(userId string) (bool, User, error) {
	var exists bool
	var user User

	// Query to check if the user exists based on the UserId
	query := `SELECT UserJson FROM users WHERE UserId = ?`
	row := db.QueryRow(query, userId)

	// Retrieve the JSON data from the row
	var userJson []byte
	err := row.Scan(&userJson)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			return false, user, nil
		}
		return false, user, fmt.Errorf("error checking user existence: %v", err)
	}

	// Deserialize the JSON data into the User struct
	user, err = JSONToObject(userJson)
	if err != nil {
		return false, user, fmt.Errorf("error unmarshalling user JSON: %v", err)
	}

	// If the user was found and deserialized, it exists
	exists = true
	return exists, user, nil
}

func DeleteUser(userId string) error {
	query := `DELETE FROM users WHERE UserId = ?`
	_, err := db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	fmt.Println("SQLite database connection closed.")
}
func PrintAllUsers() error {
	// Define a query to fetch all users
	query := `SELECT UserId, UserJson FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error fetching all users: %v", err)
	}
	defer rows.Close() // Ensure the rows are closed when done

	// Iterate over the rows and print each user's information
	for rows.Next() {
		var userId string
		var userJson []byte
		// Scan each row into the variables
		err := rows.Scan(&userId, &userJson)
		if err != nil {
			return fmt.Errorf("error scanning user data: %v", err)
		}

		// Unmarshal the JSON data into the User struct
		var user User
		err = json.Unmarshal(userJson, &user)
		if err != nil {
			return fmt.Errorf("error unmarshalling user data: %v", err)
		}

		// Print user information
		fmt.Printf("UserID: %s\n", user.UserId)
		fmt.Printf("UserName: %s\n", user.UserName)
		fmt.Printf("Points: %d\n", user.Points)
		fmt.Printf("JoinedDate: %s\n", user.JoinedDate)
		fmt.Printf("FirstVideoLink: %s\n", user.FirstVideoLink)
		fmt.Printf("LastVideoLink: %s\n", user.LastVideoLink)
		fmt.Printf("LastComment: %s\n", user.LastComment)
		fmt.Printf("LastSeen: %s\n", user.LastSeen)
		fmt.Printf("ProfilePic: %s\n", user.ProfilePic)
		fmt.Println("-------------------------------")
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %v", err)
	}

	return nil
}

func AddColumnToUsers(columnName string, columnType string) error {
	// Define the SQL statement to add a new column to the 'users' table
	query := fmt.Sprintf(`ALTER TABLE users ADD COLUMN %s %s;`, columnName, columnType)

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error adding column: %v", err)
	}

	fmt.Printf("Successfully added column: %s %s\n", columnName, columnType)
	return nil
}
