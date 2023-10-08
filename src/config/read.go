package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	dbHost     = "DB_HOST"
	dbUser     = "DB_USER"
	dbName     = "DB_NAME"
	dbPort     = "DB_PORT"
	dbPassword = "DB_PASSWORD"
)

type configVariables struct {
	dbHost     string
	dbUser     string
	dbName     string
	dbPort     string
	dbPassword string
}

func Read() (res Config, err error) {
	var cfg configVariables
	var envMap = map[string]*string{
		dbHost:     &cfg.dbHost,
		dbUser:     &cfg.dbUser,
		dbName:     &cfg.dbName,
		dbPort:     &cfg.dbPort,
		dbPassword: &cfg.dbPassword,
	}

	defer func() {
		for envKey, _ := range envMap {
			unsetErr := os.Unsetenv(envKey)
			if unsetErr != nil {
				err = fmt.Errorf("internal err: %w, unset env: %v", err, unsetErr)
				break
			}
		}
	}()

	for envKey, cfgValue := range envMap {
		val, ok := os.LookupEnv(envKey)
		if !ok {
			return Config{}, fmt.Errorf("no env %s provided", envKey)
		}

		*cfgValue = val
	}

	return formatConfig(cfg)
}

func formatConfig(cfg configVariables) (Config, error) {
	port, err := strconv.Atoi(cfg.dbPort)
	if err != nil {
		return Config{}, fmt.Errorf("port in not a number: port=%s", cfg.dbPort)
	}

	return Config{
		dbHost:     cfg.dbHost,
		dbUser:     cfg.dbUser,
		dbName:     cfg.dbName,
		dbPort:     uint16(port),
		dbPassword: cfg.dbPassword,
	}, nil
}
