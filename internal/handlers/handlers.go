package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/WOLFraMix/redis-counter/internal/storage"
)

type API struct {
	Store storage.Counter
}

// writeJSON формирует HTTP‑ответ в формате JSON
func writeJSON(w http.ResponseWriter, status int, v any) {
	// устанавливает заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	// отправляет HTTP‑статус клиенту
	w.WriteHeader(status)
	// создаёт JSON‑энкодер,
	// записывающий данные напрямую в http.ResponseWriter,
	// и сериализует значение v в JSON‑формат
	json.NewEncoder(w).Encode(v)
}

// POST /hit/{name}
func (a *API) Hit(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	n, err := a.Store.Hit(r.Context(), name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"name": name, "hits": n})
}

// GET /stats/{name}
func (a *API) Stats(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	n, err := a.Store.Stats(r.Context(), name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"name": name, "hits": n})
}

// GET /top?limit=10
func (a *API) Top(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10 // дефолт и защита от «top?limit=999999»
	}
	names, err := a.Store.Top(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"top": names})
}
