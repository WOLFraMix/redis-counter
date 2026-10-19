# redis-counter

Счётчик событий и топ популярных имён на Go + Redis.

## Запуск

1. docker compose up --build
2. Сервер: <http://localhost:8080> (Redis уже в сети compose)

## API

- POST /hit/{name} — засчитать хит → {"name": "...", "hits": N}
- GET /stats/{name} — сколько всего хитов (несуществующее имя → 0)
- GET /top?limit=10 — топ имён по числу хитов
- GET /healthz — проверка живости

## Тесты

go test ./...
