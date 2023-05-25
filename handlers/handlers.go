package handlers

import (
	"bookswap-service-tdd-go/db"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Handler struct {
	bs *db.BookService
	us *db.UserService
}

func NewHandler(bs *db.BookService, us *db.UserService) *Handler {
	return &Handler{
		bs: bs,
		us: us,
	}
}

func (h *Handler) Index(c *fiber.Ctx) error {
	resp := &Response{
		Message: "Welcome to the BookSwap service!",
		Books:   h.bs.List(),
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *Handler) ListBooks(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(&Response{
		Books: h.bs.List(),
	})
}

func (h *Handler) UserUpsert(c *fiber.Ctx) error {
	var user db.User
	err := c.BodyParser(&user)
	if err != nil {
		errorMessage := fmt.Errorf("invalid user body:%v", err)
		return c.Status(http.StatusBadRequest).
			JSON(&Response{Error: errorMessage.Error()})
	}

	user, err = h.us.Upsert(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(&Response{User: &user})
}

func (h *Handler) ListUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	user, books, err := h.us.Get(userID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&Response{Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		User:  user,
		Books: books,
	})
}

func (h *Handler) SwapBook(c *fiber.Ctx) error {
	bookID := c.Params("id")
	userID := c.Query("user")

	err := h.us.Exists(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{Error: err.Error()})
	}

	_, err = h.bs.SwapBook(bookID, userID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&Response{Error: err.Error()})
	}

	user, books, err := h.us.Get(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		User:  user,
		Books: books,
	})
}

func (h *Handler) BookUpsert(c *fiber.Ctx) error {
	var book db.Book
	err := c.BodyParser(&book)
	if err != nil {
		invalidBodyError := fmt.Errorf("invalid book body:%v", err)
		return c.Status(http.StatusInternalServerError).JSON(&Response{Error: invalidBodyError.Error()})
	}

	err = h.us.Exists(book.OwnerID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{Error: err.Error()})
	}

	book = h.bs.Upsert(book)
	return c.Status(http.StatusOK).JSON(&Response{Books: []db.Book{book}})
}
