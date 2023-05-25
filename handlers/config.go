package handlers

import "github.com/gofiber/fiber/v2"

func ConfigureServer(fiberApp *fiber.App, handler *Handler) {
	fiberApp.Get("/", handler.Index)
	fiberApp.Get("/books", handler.ListBooks)
	fiberApp.Post("/users", handler.UserUpsert)
	fiberApp.Get("/users/{id}", handler.ListUserByID)
	fiberApp.Post("/books/{id}", handler.SwapBook)
	fiberApp.Post("/books", handler.BookUpsert)
}
