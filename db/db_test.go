package db_test

import (
	"errors"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/mattn/go-sqlite3"
)

func TestNewDB(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path string
		want sqlite3.ErrNoExtended
	}{
		"Normal":                {path: "../.sqlite3/db_test.db"},
		"Directory not created": {path: "../.nothing/db_test.db", want: sqlite3.ErrCantOpenFullPath},
	}

	t.Cleanup(func() {
		if err := os.Remove(cases["Normal"].path); err != nil && !os.IsNotExist(err) {
			t.Error("failed to cleanup testdata, err =", err)
		}
	})

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := db.NewDB(c.path)
			if c.want == 0 { // エラーを期待しないケース
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				defer dbConn.Close()
				return
			}

			var sqliteErr sqlite3.Error
			if !errors.As(err, &sqliteErr) || sqliteErr.ExtendedCode != c.want {
				t.Errorf("unexpected error, got = %v, want sqlite extended error code %d", err, c.want)
				return
			}
		})
	}
}
