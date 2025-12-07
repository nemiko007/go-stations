package sta19_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func TestStation19(t *testing.T) {
	testcases := map[string]struct {
		IDs                []string
		InitialDataCount   int
		WantHTTPStatusCode int
	}{
		"Empty Ids": {
			IDs:                []string{},
			InitialDataCount:   3,
			WantHTTPStatusCode: http.StatusOK,
		},
		"Not found ID": {
			IDs:                []string{"4"},
			InitialDataCount:   3,
			WantHTTPStatusCode: http.StatusOK,
		},
		"One delete": {
			IDs:                []string{"1"},
			InitialDataCount:   3,
			WantHTTPStatusCode: http.StatusOK,
		},
		"Multiple delete": {
			IDs:                []string{"2", "3"},
			InitialDataCount:   3,
			WantHTTPStatusCode: http.StatusOK,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			dbPath := fmt.Sprintf("./%s.db", strings.ReplaceAll(t.Name(), "/", "_"))
			if err := os.Setenv("DB_PATH", dbPath); err != nil {
				t.Fatalf("dbPathのセットに失敗しました。%v", err)
			}

			todoDB, err := db.NewDB(dbPath)
			if err != nil {
				t.Fatalf("データベースの作成に失敗しました: %v", err)
			}

			t.Cleanup(func() {
				if err := todoDB.Close(); err != nil {
					t.Errorf("データベースのクローズに失敗しました: %v", err)
				}
				if err := os.Remove(dbPath); err != nil {
					t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
				}
			})

			stmt, err := todoDB.Prepare(`INSERT INTO todos(subject) VALUES(?)`)
			if err != nil {
				t.Fatalf("ステートメントの作成に失敗しました: %v", err)
			}
			defer stmt.Close()

			for i := 0; i < tc.InitialDataCount; i++ {
				if _, err := stmt.Exec("subject"); err != nil {
					t.Fatalf("todoの追加に失敗しました: %v", err)
				}
			}

			r := router.NewRouter(todoDB)
			srv := httptest.NewServer(r)
			defer srv.Close()

			req, err := http.NewRequest(http.MethodDelete, srv.URL+"/todos",
				bytes.NewBufferString(fmt.Sprintf(`{"ids":[%s]}`, strings.Join(tc.IDs, ","))))
			if err != nil {
				t.Errorf("リクエストの作成に失敗しました: %v", err)
				return
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("リクエストの送信に失敗しました: %v", err)
				return
			}
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("レスポンスのクローズに失敗しました: %v", err)
					return
				}
			})

			if resp.StatusCode != tc.WantHTTPStatusCode {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, tc.WantHTTPStatusCode)
				return
			}
		})
	}
}
