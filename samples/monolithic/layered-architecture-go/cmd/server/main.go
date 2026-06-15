// Command server は各層を結線（DI）して HTTP サーバを起動するエントリポイント。
package main

import (
	"log"
	"net/http"
	"os"

	"example.com/layered-architecture-go/internal/application"
	"example.com/layered-architecture-go/internal/infrastructure"
	"example.com/layered-architecture-go/internal/presentation"
)

func main() {
	// 依存の組み立て: インフラ → アプリ → プレゼンテーション。
	repo := infrastructure.NewMemoryTaskRepository()
	svc := application.NewTaskService(repo)
	handler := presentation.NewHandler(svc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler.Routes()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
