package config

import (
	"fmt"
	"os"
)

const (
	dbHost     = "DB_HOST"
	dbUser     = "DB_USER"
	dbName     = "DB_NAME"
	dbPort     = "DB_PORT"
	dbPassword = "DB_PASSWORD"
)

type Config struct {
	dbHost     string
	dbUser     string
	dbName     string
	dbPort     string
	dbPassword string
}

func Read() (res Config, err error) {
	var envMap = map[string]*string{
		dbHost:     &res.dbHost,
		dbUser:     &res.dbUser,
		dbName:     &res.dbName,
		dbPort:     &res.dbPort,
		dbPassword: &res.dbPassword,
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

	return res, nil
}

func (c Config) DbHost() string {
	return c.dbHost
}

func (c Config) DbUser() string {
	return c.dbUser
}

func (c Config) DbName() string {
	return c.dbName
}

func (c Config) DbPort() string {
	return c.dbPort
}

func (c Config) DbPassword() string {
	return c.dbPassword
}
