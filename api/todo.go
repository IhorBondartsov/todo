package api

import "time"

type ToDo struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int64     `json:"user_id"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
