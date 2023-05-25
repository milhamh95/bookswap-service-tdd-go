package db

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Book struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	OwnerID string `json:"owner_id"`
	Status  string `json:"status"`
}

type BookService struct {
	books map[string]Book
	ps    PostingService
}

func NewBookService(initial []Book, ps PostingService) *BookService {
	books := make(map[string]Book)
	for _, b := range initial {
		books[b.ID] = b
	}
	return &BookService{
		books: books,
		ps:    ps,
	}
}

func (bs *BookService) Get(id string) (*Book, error) {
	b, ok := bs.books[id]
	if !ok {
		return nil, errors.New("no book found")
	}

	return &b, nil
}

func (bs *BookService) Upsert(b Book) Book {
	_, ok := bs.books[b.ID]
	if !ok {
		b.ID = uuid.NewString()
		b.Status = Available.String()
	}
	bs.books[b.ID] = b

	return b
}

func (bs *BookService) List() []Book {
	var items []Book
	for _, b := range bs.books {
		if b.Status == Available.String() {
			items = append(items, b)
		}
	}

	return items
}

func (bs *BookService) ListByUser(userID string) []Book {
	var items []Book
	for _, b := range bs.books {
		if b.OwnerID == userID {
			items = append(items, b)
		}
	}
	return items
}

func (bs *BookService) SwapBook(bookID, userID string) (*Book, error) {
	book, ok := bs.books[bookID]
	if !ok {
		return nil, fmt.Errorf("no book found for id %s", bookID)
	}
	if book.Status != Available.String() {
		return nil, fmt.Errorf("book %s is not available for swapping", bookID)
	}
	book.OwnerID = userID
	book.Status = Swapped.String()
	bs.books[bookID] = book
	return &book, nil
}
