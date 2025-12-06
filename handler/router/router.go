package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

// NewRouter はエンドポイントを登録して http.Handler を返す
func NewRouter(todoDB *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// 例: /health にアクセスすると "ok" を返す
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)
	// 例: /todos にアクセスすると TodoHandler が処理する
	mux.HandleFunc("/todos/", todoHandler.ServeHTTP)

	return mux
}
