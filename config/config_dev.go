//go:build dev

package config

const (
	REDIS_DATABASE = "localhost:6379"
	DB_HOST        = "./db/dev.db"
	DB_DRIVER      = "sqlite3"
	API_PORT       = 8080
)
