package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	port        string
	environment string
	logLevel    string

	accessTokenTimeout  string
	refreshTokenTimeout string
	jwtSecret           string

	postgresHost     string
	postgresPort     string
	postgresUser     string
	postgresPassword string
	postgresDb       string
	postgresSslmode  string
}

func (c *Config) Port() string        { return c.port }
func (c *Config) Environment() string { return c.environment }
func (c *Config) LogLevel() string    { return c.logLevel }

func (c *Config) AccessTokenTimeout() string  { return c.accessTokenTimeout }
func (c *Config) RefreshTokenTimeout() string { return c.refreshTokenTimeout }
func (c *Config) JWTSecret() string           { return c.jwtSecret }

func (c *Config) PostgresHost() string     { return c.postgresHost }
func (c *Config) PostgresPort() string     { return c.postgresPort }
func (c *Config) PostgresUser() string     { return c.postgresUser }
func (c *Config) PostgresPassword() string { return c.postgresPassword }
func (c *Config) PostgresDb() string       { return c.postgresDb }
func (c *Config) PostgresSslmode() string  { return c.postgresSslmode }

var globalCfg Config

func Load() error {
	globalCfg = Defaults

	godotenv.Load() // optional .env, ignore if missing

	if v := os.Getenv("APP_PORT"); v != "" {
		globalCfg.port = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		globalCfg.logLevel = v
	}

	if v := os.Getenv("ACCESS_TOKEN_TIMEOUT"); v != "" {
		globalCfg.accessTokenTimeout = v
	}
	if v := os.Getenv("REFRESH_TOKEN_TIMEOUT"); v != "" {
		globalCfg.refreshTokenTimeout = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		globalCfg.jwtSecret = v
	}

	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		globalCfg.postgresHost = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		globalCfg.postgresPort = v
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		globalCfg.postgresUser = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		globalCfg.postgresPassword = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		globalCfg.postgresDb = v
	}
	if v := os.Getenv("POSTGRES_SSLMODE"); v != "" {
		globalCfg.postgresSslmode = v
	}

	return nil
}

func Get() Config {
	return globalCfg
}
