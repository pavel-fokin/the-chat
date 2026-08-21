package thechat

import "github.com/google/uuid"

type user struct{}

// UserID identifies a User.
type UserID = ID[user]

// NewUserID returns a new UserID.
func NewUserID() UserID {
	return ID[user]{uuid.Must(uuid.NewV7())}
}

// User is a person using the chat.
type User struct {
	ID UserID
}
