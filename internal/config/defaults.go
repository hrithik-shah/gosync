package config

var Defaults = Config{
	port:        "8080",
	environment: "development",
	logLevel:    "info",

	accessTokenTimeout:  "15s",
	refreshTokenTimeout: "14d",

	postgresHost:    "localhost",
	postgresPort:    "5432",
	postgresUser:    "postgres",
	postgresDb:      "app",
	postgresSslmode: "disable",
}
