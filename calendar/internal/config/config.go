package config

// Application holds application configuration values
type Application struct {
	DB *Database
}

type Database struct {
	DSN string `env:"DSN"`
}
