package config

import "os"

type Config struct {
	Port      string
	RedisAddr string
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // удобно для запуска без докера
	}
	return &Config{Port: port, RedisAddr: addr}
}
