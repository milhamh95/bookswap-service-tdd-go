package db

import (
	"fmt"
	"github.com/rs/xid"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
}

type BookOperationService interface {
	ListByUser(userID string) []Book
}

type UserService struct {
	users map[string]User
	bs    BookOperationService
}

func NewUserService(initial []User, bs BookOperationService) *UserService {
	users := make(map[string]User)
	for _, u := range initial {
		users[u.ID] = u
	}
	return &UserService{
		users: users,
		bs:    bs,
	}
}

// Get returns a given user or error if none exists.
func (us *UserService) Get(id string) (*User, []Book, error) {
	u, ok := us.users[id]
	if !ok {
		return nil, nil, fmt.Errorf("no user found for id %s", id)
	}
	books := us.bs.ListByUser(id)

	return &u, books, nil
}

// Exists returns whether a given user exists and returns an error if none found.
func (us *UserService) Exists(id string) error {
	_, ok := us.users[id]
	if !ok {
		return fmt.Errorf("no user found for id %s", id)
	}

	return nil
}

// Upsert creates or updates a new order.
func (us *UserService) Upsert(u User) (User, error) {
	u.ID = xid.New().String()
	us.users[u.ID] = u

	return u, nil
}
