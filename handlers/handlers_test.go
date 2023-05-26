package handlers_test

import (
	"bookswap-service-tdd-go/db"
	"bookswap-service-tdd-go/handlers"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIndexIntegration(t *testing.T) {
	t.Run("with fiber", func(t *testing.T) {
		if os.Getenv("LONG") == "" {
			t.Skip("Skipping TestIndexIntegration in short mode")
		}

		// Arrange
		book := db.Book{
			ID:     xid.New().String(),
			Name:   "my first integration test",
			Status: db.Available.String(),
		}
		bs := db.NewBookService([]db.Book{book}, nil)
		ha := handlers.NewHandler(bs, nil)

		app := fiber.New()
		handlers.ConfigureServer(app, ha)

		req, _ := http.NewRequest(
			"GET",
			"/",
			nil,
		)

		res, err := app.Test(req, -1)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		require.NoError(t, err)

		var resp handlers.Response
		err = json.Unmarshal(body, &resp)
		require.NoError(t, err)

		require.Equal(t, 1, len(resp.Books))
		require.Contains(t, resp.Books, book)
	})

	t.Run("with http test new server", func(t *testing.T) {
		// Arrange
		book := db.Book{
			ID:     xid.New().String(),
			Name:   "my first integration test",
			Status: db.Available.String(),
		}
		bs := db.NewBookService([]db.Book{book}, nil)
		ha := handlers.NewHandler(bs, nil)

		svr := httptest.NewServer(http.HandlerFunc(adaptor.FiberHandlerFunc(ha.Index)))
		defer svr.Close()

		r, err := http.Get(svr.URL)
		require.Equal(t, http.StatusOK, r.StatusCode)
		require.NoError(t, err)

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		require.NoError(t, err)

		var resp handlers.Response
		err = json.Unmarshal(body, &resp)
		require.NoError(t, err)

		require.Equal(t, 1, len(resp.Books))
		require.Contains(t, resp.Books, book)
	})
}
