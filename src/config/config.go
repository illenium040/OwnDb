package config

type Config struct {
	dbHost     string
	dbUser     string
	dbName     string
	dbPort     uint16
	dbPassword string
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

func (c Config) DbPort() uint16 {
	return c.dbPort
}

func (c Config) DbPassword() string {
	return c.dbPassword
}
