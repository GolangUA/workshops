//go:build integration

package postgre

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const userName = "user"

type TestSuite struct {
	suite.Suite
	repo *Repository
	user *models.User
}

func (s *TestSuite) cleanup(silent bool) {
	var err error
	_, err = psql.Delete("user_event").RunWith(s.repo.db).Exec()
	if !silent {
		require.NoError(s.T(), err, "cleanup user_event")
	}
	_, err = psql.Delete("event").RunWith(s.repo.db).Exec()
	if !silent {
		require.NoError(s.T(), err, "cleanup event")
	}
	_, err = psql.Delete("users").RunWith(s.repo.db).Exec()
	if !silent {
		require.NoError(s.T(), err, "cleanup user")
	}
}

func (s *TestSuite) SetupSuite() {
	s.repo = NewRepository("host=127.0.0.1 port=5432 user=goit password=goit dbname=goit sslmode=disable")
}

func (s *TestSuite) SetupTest() {
	s.cleanup(false)
	var err error
	s.user, err = s.repo.CreateUser(userName, "password", "UTC")
	require.NoError(s.T(), err, "create user")
}

func (s *TestSuite) TearDownSuite() {
	s.cleanup(true)
	_ = s.repo.db.Close()
}

func (s *TestSuite) dateTime(val string) time.Time {
	res, err := time.ParseInLocation("2006-01-02 15:04", val, time.Local)
	require.NoErrorf(s.T(), err, "parse %s", val)
	return res
}

func Test_PostgreTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
