package db

type PostingService interface {
	NewOrder(b Book) error
}
