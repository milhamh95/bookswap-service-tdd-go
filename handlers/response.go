package handlers

import "bookswap-service-tdd-go/db"

type Response struct {
	Message string    `json:"message,omitempty"`
	Error   string    `json:"error,omitempty"`
	Books   []db.Book `json:"books,omitempty"`
	User    *db.User  `json:"user,omitempty"`
}
