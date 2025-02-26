package config

import (
	"fmt"
	"github.com/friendsofgo/errors"
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

	switch {
	case c.Postgres.User == "":
		return errors.New("Postgres.User is empty")
	case c.Postgres.Password == "":
		return errors.New("Postgres.Password is empty")
	case c.Postgres.Host == "":
		return errors.New("Postgres.Host is empty")
	case c.Postgres.Port == "":
		return errors.New("Postgres.Port is empty")
	case c.Postgres.DB == "":
		return errors.New("Postgres.DB is empty")
	case c.ServerAddress == "":
		return errors.New("ServerAddress is empty")
	}
	return nil
}

func (c *PostgresConfig) PostgresConnection() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DB)
}
