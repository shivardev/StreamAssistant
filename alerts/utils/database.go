package utils

import (
	"encoding/json"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

var db *badger.DB

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
	// Open a Badger database
	var err error
	db, err = badger.Open(badger.DefaultOptions("./badgerDB"))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

}

func GetUserById(userId string) (User, error) {
	var user User
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userId))
		if err != nil {
			return err
		}
		fmt.Println("Item:", item)
		return nil
	})
	if err != nil {
		return user, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

// Insert or update a user in the database
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
	data, err := ObjectToJSON(&insertUser)
	if err != nil {
		return insertUser, fmt.Errorf("error converting user to JSON: %v", err)
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(msg.AuthorId), []byte(data))
		return err
	})
	if err != nil {
		return insertUser, fmt.Errorf("error inserting or updating user: %v", err)
	}
	return insertUser, nil
}

func UpdateUser(msg ChatMessage, user User) (User, error) {
	var updatedUser = User{
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
	data, err := ObjectToJSON(&updatedUser)
	if err != nil {
		return updatedUser, fmt.Errorf("error converting user to JSON: %v", err)
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(msg.AuthorId), []byte(data))
		return err
	})
	if err != nil {
		return updatedUser, fmt.Errorf("error inserting or updating user: %v", err)
	}
	return user, nil
}

// Check if a user exists by UserId
func CheckUserExists(userId string) (bool, User, error) {
	var exists bool
	var valCopy []byte
	var user User

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userId))
		if err == nil {
			exists = true
		} else if err == badger.ErrKeyNotFound {
			exists = false
		} else {
			return fmt.Errorf("error getting user from db: %v", err)
		}
		if item == nil {
			exists = false
			return nil
		}
		err2 := item.Value(func(val []byte) error {
			fmt.Printf("The UserData is: %s\n", val)
			valCopy = append([]byte{}, val...)
			return nil
		})

		if err2 != nil {
			return err2
		}

		return nil
	})

	if err != nil {
		return false, user, fmt.Errorf("error checking user existence: %v", err)
	}

	if exists {
		err := json.Unmarshal(valCopy, &user)
		if err != nil {
			return false, user, fmt.Errorf("error unmarshalling user data: %v", err)
		}
	}

	return exists, user, nil
}

// Delete a user by UserId
func DeleteUser(userId string) error {
	err := db.Update(func(txn *badger.Txn) error {
		// Delete the user by UserId
		err := txn.Delete([]byte(userId))
		return err
	})

	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// Close the database connection
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	fmt.Println("Badger database connection closed.")
}
