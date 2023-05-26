package handlers_test

import (
	"bookswap-service-tdd-go/db"
	"bookswap-service-tdd-go/handlers"
	"encoding/json"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/xid"
	"io"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Handlers integration", func() {
	var svr *httptest.Server
	var eb db.Book

	BeforeEach(func() {
		eb = db.Book{
			ID:     xid.New().String(),
			Name:   "first integration test",
			Status: db.Available.String(),
		}
		bs := db.NewBookService([]db.Book{eb}, nil)
		ha := handlers.NewHandler(bs, nil)
		svr = httptest.NewServer(http.Handler(adaptor.FiberHandlerFunc(ha.Index)))
	})

	AfterEach(func() {
		svr.Close()
	})

	Describe("index endpoint", func() {
		Context("with one existing book", func() {
			It("should return book", func() {
				r, err := http.Get(svr.URL)
				Expect(err).To(BeNil())
				Expect(r.StatusCode).To(Equal(http.StatusOK))

				body, err := io.ReadAll(r.Body)
				r.Body.Close()
				Expect(err).To(BeNil())

				var resp handlers.Response
				err = json.Unmarshal(body, &resp)

				Expect(err).To(BeNil())
				Expect(len(resp.Books)).To(Equal(1))
				Expect(resp.Books).To(ContainElement(eb))
			})
		})
	})
})
