package main

import (
	"context"
	"log"
	"net/http"

	"github.com/WOLFraMix/redis-counter/internal/config"
	"github.com/WOLFraMix/redis-counter/internal/storage"
)

func main() {
	cfg := config.New()

	st, err := storage.New(context.Background(), cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis connect: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	_ = st // заглушка

	log.Printf("server listening on %s (redis: %s)", cfg.Port, cfg.RedisAddr)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
