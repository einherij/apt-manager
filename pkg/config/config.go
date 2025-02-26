package config

import (
	"fmt"
	"os"
)

type Config struct {
	Postgres      PostgresConfig
	ServerAddress string
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

func NewConfig() *Config {
	return new(Config)
}

func (c *Config) ParseEnv() error {
	c.Postgres.User = os.Getenv("PG_USER")
	c.Postgres.Password = os.Getenv("PG_PASSWORD")
	c.Postgres.Host = os.Getenv("PG_HOST")
	c.Postgres.Port = os.Getenv("PG_PORT")
	c.Postgres.DB = os.Getenv("PG_DB")

	c.ServerAddress = os.Getenv("SERVER_ADDRESS")

	if c.Postgres.User == "" ||
		c.Postgres.Password == "" ||
		c.Postgres.Host == "" ||
		c.Postgres.Port == "" ||
		c.Postgres.DB == "" ||
		c.ServerAddress == "" {
		return fmt.Errorf("config field is empty: %v", c)
	}
	return nil
}

func (c *PostgresConfig) PostgresConnection() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DB)
}
