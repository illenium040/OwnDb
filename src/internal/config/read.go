package config

import (
	"fmt"
	"os"
)

const (
	dbUrl = "DB_URL"
)

func Read() (cfg Config, err error) {
	var envMap = map[string]*string{
		dbUrl: &cfg.dbUrl,
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

	return cfg, nil
}
