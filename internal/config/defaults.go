package config

import "time"

var Defaults = Config{
	port:        "8080",
	environment: "development",
	logLevel:    "info",

	accessTokenTimeout:  time.Minute * 15,    // 15m
	refreshTokenTimeout: time.Hour * 24 * 14, // 14d

	postgresHost:    "localhost",
	postgresPort:    "5432",
	postgresUser:    "postgres",
	postgresDb:      "app",
	postgresSslmode: "disable",
}
