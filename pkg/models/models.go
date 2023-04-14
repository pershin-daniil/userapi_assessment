package models

import "time"

type (
	UserList  map[string]User
	UserStore struct {
		Increment int      `json:"increment"`
		List      UserList `json:"list"`
	}
)

type UserRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type User struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
