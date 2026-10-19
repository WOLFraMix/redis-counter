package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WOLFraMix/redis-counter/internal/storage"
)

func testAPI() *API {
	return &API{Store: storage.NewFake()}
}

func serve(t *testing.T, method, target string, api *API) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /hit/{name}", api.Hit)
	mux.HandleFunc("GET /stats/{name}", api.Stats)
	mux.HandleFunc("GET /top", api.Top)

	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHitIncrements(t *testing.T) {
	api := testAPI()

	serve(t, "POST", "/hit/go-dev", api)
	serve(t, "POST", "/hit/go-dev", api)
	rec := serve(t, "GET", "/stats/go-dev", api)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]any
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body["hits"].(float64) != 2 {
		t.Fatalf("hits = %v, want 2", body["hits"])
	}
}

func TestStatsUnknown(t *testing.T) {
	api := testAPI()
	rec := serve(t, "GET", "/stats/nope", api)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]any
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body["hits"].(float64) != 0 {
		t.Fatalf("hits = %v, want 0", body["hits"])
	}
}

func TestTopOrder(t *testing.T) {
	api := testAPI()
	// news накрутили больше всех, go-dev один раз
	for i := 0; i < 3; i++ {
		serve(t, "POST", "/hit/news", api)
	}
	serve(t, "POST", "/hit/go-dev", api)

	rec := serve(t, "GET", "/top", api)
	var body map[string][]string
	json.Unmarshal(rec.Body.Bytes(), &body)

	got := body["top"]
	if len(got) != 2 || got[0] != "news" || got[1] != "go-dev" {
		t.Fatalf("top = %v, want [news go-dev]", got)
	}
}
