package handlers_test

import (
	"bookswap-service-tdd-go/db"
	"bookswap-service-tdd-go/handlers"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestIndexIntegration(t *testing.T) {
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
}

func TestIndexRoute(t *testing.T) {
	// Define a structure for specifying input and output
	// data of a single test case. This structure is then used
	// to create a so called test map, which contains all test
	// cases, that should be run for testing this function
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "index route",
			route:         "/",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "OK",
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "Cannot GET /i-dont-exist",
		},
	}

	// Setup the app as it is done in the main function
	// Initialize a new app
	app := fiber.New()

	// Register the index route with a simple
	// "OK" response. It should return status
	// code 200
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
