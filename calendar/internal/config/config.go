package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

// Application holds application configuration values
type Application struct {
	DB     *Database `yaml:"db"`
	BCrypt *BCrypt   `yaml:"bcrypt"`
	JWT    *JWT      `json:"jwt"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
}

type BCrypt struct {
	Secret string `yaml:"secret"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

var instance *Application
var once sync.Once

func GetConfig() *Application {
	once.Do(initApplicationConfig)
	return instance
}

func initApplicationConfig() {
	instance = &Application{
		DB:     &Database{},
		BCrypt: &BCrypt{},
		JWT:    &JWT{},
	}

	if err := yaml.Unmarshal(configBytes(), instance); err != nil {
		panic(err)
	}

	overrideByEnv(&instance.DB.Host, "DB_HOST", "localhost")
	overrideByEnv(&instance.DB.Port, "DB_PORT", "5432")
	overrideByEnv(&instance.DB.Name, "DB_NAME", "gotest")
	overrideByEnv(&instance.DB.User, "DB_USER", "gouser")
	overrideByEnv(&instance.DB.Password, "DB_PASSWORD", "gopassword")
	overrideByEnv(&instance.DB.SSLMode, "DB_SSL_MODE", "")
	overrideByEnv(&instance.BCrypt.Secret, "BCRYPT_SECRET", "")
	overrideByEnv(&instance.JWT.Secret, "JWT_SECRET", "")

	if instance.BCrypt.Secret == "" {
		panic("bcrypt.secret must not be empty")
	}
	if instance.JWT.Secret == "" {
		panic("jwt.secret must not be empty")
	}
}

func (a *Application) DSN() string {
	res := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", a.DB.Host, a.DB.Port, a.DB.User, a.DB.Password, a.DB.Name)
	if a.DB.SSLMode != "" {
		res = res + " sslmode=" + a.DB.SSLMode
	}
	return res
}

var configBytes = func() []byte {
	res, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	return res
}

func overrideByEnv(target *string, key string, def string) {
	if value := os.Getenv(key); value != "" {
		*target = value
	} else if *target == "" && def != "" {
		*target = def
	}
}
