package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestApplication_DSN_defaultValues(t *testing.T) {
	original := configBytes
	defer func() {
		configBytes = original
	}()
	configBytes = func() []byte {
		return []byte(`bcrypt:
  secret: secret
jwt:
  secret: secret`)
	}

	initApplicationConfig()
	assert.Equal(t, "host=localhost port=5432 user=gouser password=gopassword dbname=gotest", GetConfig().DSN())
}

func TestApplication_DSN_parsed(t *testing.T) {
	original := configBytes
	defer func() {
		configBytes = original
	}()
	configBytes = func() []byte {
		return []byte((`db:
  host: 127.0.0.1
  port: 2345
  name: db
  user: user
  password: password
  sslmode: other
bcrypt:
  secret: bcrypt_secret
jwt:
  secret: jwt_secret`))
	}

	initApplicationConfig()
	assert.Equal(t, "host=127.0.0.1 port=2345 user=user password=password dbname=db sslmode=other", GetConfig().DSN())
}

func TestApplication_DSN_overridable(t *testing.T) {
	original := configBytes
	defer func() {
		configBytes = original
	}()
	configBytes = func() []byte {
		return []byte(`bcrypt:
  secret: bcrypt_secret
jwt:
  secret: jwt_secret`)
	}

	type override struct {
		env      string
		value    string
		expected string
	}

	overrides := []override{
		{"DB_HOST", "overridden", "host=overridden port=5432 user=gouser password=gopassword dbname=gotest"},
		{"DB_PORT", "overridden", "host=localhost port=overridden user=gouser password=gopassword dbname=gotest"},
		{"DB_USER", "overridden", "host=localhost port=5432 user=overridden password=gopassword dbname=gotest"},
		{"DB_PASSWORD", "overridden", "host=localhost port=5432 user=gouser password=overridden dbname=gotest"},
		{"DB_NAME", "overridden", "host=localhost port=5432 user=gouser password=gopassword dbname=overridden"},
		{"DB_SSL_MODE", "disable", "host=localhost port=5432 user=gouser password=gopassword dbname=gotest sslmode=disable"},
	}

	for _, o := range overrides {
		os.Setenv(o.env, o.value)
		initApplicationConfig()
		assert.Equalf(t, o.expected, instance.DSN(), "%s does not override application config", o.env)
		assert.Equal(t, "bcrypt_secret", instance.BCrypt.Secret)
		assert.Equal(t, "jwt_secret", instance.JWT.Secret)
		os.Unsetenv(o.env)
	}

	os.Setenv("BCRYPT_SECRET", "overridden")
	initApplicationConfig()
	assert.Equalf(t, "overridden", instance.BCrypt.Secret, "BCRYPT_SECRET does not override application config")
	os.Unsetenv("BCRYPT_SECRET")

	os.Setenv("JWT_SECRET", "overridden")
	initApplicationConfig()
	assert.Equalf(t, "overridden", instance.JWT.Secret, "JWT_SECRET does not override application config")
	os.Unsetenv("JWT_SECRET")
}
