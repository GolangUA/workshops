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
	configBytes = func() []byte { return []byte("") }

	initApplicationConfig()
	assert.Equal(t, "host=localhost port=5432 user=gouser password=gopassword dbname=gotest", GetConfig().DSN())
}

func TestApplication_DSN_parsed(t *testing.T) {
	original := configBytes
	defer func() {
		configBytes = original
	}()
	configBytes = func() []byte {
		return []byte(`db:
  host: 127.0.0.1
  port: 2345
  name: db
  user: user
  password: password
  sslmode: other`)
	}

	initApplicationConfig()
	assert.Equal(t, "host=127.0.0.1 port=2345 user=user password=password dbname=db sslmode=other", GetConfig().DSN())
}

func TestApplication_DSN_overridable(t *testing.T) {
	original := configBytes
	defer func() {
		configBytes = original
	}()
	configBytes = func() []byte { return []byte("") }

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
		os.Unsetenv(o.env)
	}
}
