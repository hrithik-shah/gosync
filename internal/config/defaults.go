package config

var Defaults = Config{
	Port:        "8080",
	Environment: "development",
	LogLevel:    "info",

	POSTGRES_HOST:    "localhost",
	POSTGRES_PORT:    "5432",
	POSTGRES_USER:    "postgres",
	POSTGRES_DB:      "app",
	POSTGRES_SSLMODE: "disable",
}
