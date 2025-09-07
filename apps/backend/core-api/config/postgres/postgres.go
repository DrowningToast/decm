package postgres

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	pgxslog "github.com/mcosta74/pgx-slog"
)

const (
	MinConnections  = 0
	MaxConnections  = 10
	DefaultLogLevel = tracelog.LogLevelError
)

type Config struct {
	Host     string `env:"HOST" mapstructure:"host"`         // Default is 127.0.0.1
	Port     string `env:"PORT" mapstructure:"port"`         // Default is 5432
	User     string `env:"USER" mapstructure:"user"`         // Default is empty
	Password string `env:"PASSWORD" mapstructure:"password"` // Default is empty
	DBName   string `env:"NAME" mapstructure:"dbname"`       // Default is postgres
	SSLMode  string `env:"SSLMODE" mapstructure:"sslmode"`   // Default is prefer

	MaxConns int32 `env:"MAX_CONNS" mapstructure:"max_conns"` // Default is 8
	MinConns int32 `env:"MIN_CONNS" mapstructure:"min_conns"` // Default is 0

	Debug bool `env:"DEBUG" mapstructure:"debug"`
}

func (c *Config) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func (c *Config) QueryTracer(logger slog.Logger) pgx.QueryTracer {
	loglevel := DefaultLogLevel
	if c.Debug {
		loglevel = tracelog.LogLevelTrace
	}
	return &tracelog.TraceLog{
		Logger:   pgxslog.NewLogger(&logger),
		LogLevel: loglevel,
	}
}
