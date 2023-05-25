package db_test

import (
	"bookswap-service-tdd-go/db"
	"errors"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetBook(t *testing.T) {
	t.Run("initial books", func(t *testing.T) {
		eb := db.Book{
			ID:     xid.New().String(),
			Name:   "Existing book",
			Status: db.Available.String(),
		}
		bs := db.NewBookService([]db.Book{eb}, nil)

		tests := map[string]struct {
			id      string
			want    db.Book
			wantErr error
		}{
			"existing book": {
				id:   eb.ID,
				want: eb,
			},
			"no book found": {
				id:      "not-found",
				wantErr: errors.New("no book found"),
			},
			"empty id": {
				id:      "",
				wantErr: errors.New("no book found"),
			},
		}

		for name, tc := range tests {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				b, err := bs.Get(tc.id)
				if tc.wantErr != nil {
					require.Equal(t, tc.wantErr, err)
					require.Nil(t, b)
					return
				}

				require.Nil(t, err)
				require.Equal(t, tc.want, *b)
			})
		}
	})
}
