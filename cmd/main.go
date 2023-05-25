package cmd

import (
	"bookswap-service-tdd-go/db"
	"bookswap-service-tdd-go/handlers"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

//go:embed books.json
var booksFile []byte

//go:embed users.json
var usersFile []byte

func main() {
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	if !ok {
		log.Fatal("$BOOKSWAP_PORT not found")
	}

	books, users := importInitial()
	b := db.NewBookService(books, nil)
	u := db.NewUserService(users, b)
	h := handlers.NewHandler(b, u)

	app := fiber.New()
	handlers.ConfigureServer(app, h)

	log.Printf("Listening on :%s...\n", port)
	log.Fatal(app.Listen(":3000"))
}

func importInitial() ([]db.Book, []db.User) {
	var books []db.Book
	var users []db.User

	err := json.Unmarshal(booksFile, &books)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(usersFile, &users)
	if err != nil {
		log.Fatal(err)
	}

	return books, users
}
