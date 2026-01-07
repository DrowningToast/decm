package config

import (
	"strings"

	"github.com/cockroachdb/errors"
)

type Environment string

const (
	DevelopmentEnvironment Environment = "development"
	TestingEnvironment     Environment = "testing"
	ProductionEnvironment  Environment = "production"
)

// Parse the environment
func ParseEnvironment(s string) (Environment, error) {
	switch s {
	case "development":
		return DevelopmentEnvironment, nil
	case "testing":
		return TestingEnvironment, nil
	case "production":
		return ProductionEnvironment, nil
	default:
		return "", errors.New("invalid environment")
	}
}

// Case insensitive parse
func IParseEnvironment(s string) (Environment, error) {
	return ParseEnvironment(strings.ToLower(s))
}

func IsDevelopment() bool {
	cfg := LoadConfig()
	return cfg.Env == DevelopmentEnvironment.String()
}

func IsTesting() bool {
	cfg := LoadConfig()
	return cfg.Env == TestingEnvironment.String()
}

func IsProduction() bool {
	cfg := LoadConfig()
	return cfg.Env == ProductionEnvironment.String()
}

func (e Environment) String() string {
	return string(e)
}
