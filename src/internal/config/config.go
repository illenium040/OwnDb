package config

type Config struct {
	dbUrl string
}

func (c Config) DbUrl() string {
	return c.dbUrl
}
