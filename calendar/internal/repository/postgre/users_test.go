//go:build integration

package postgre

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) Test_GetUser() {
	user, err := s.repo.GetUser(userName)
	assert.NoError(s.T(), err, "get user")
	if assert.NotNil(s.T(), user) {
		assert.Equal(s.T(), userName, user.Name)
	}

	created, err := s.repo.CreateUser("test", "pass", "zone")
	if assert.NotNil(s.T(), created) {
		actual, err := s.repo.GetUser("test")
		assert.NoError(s.T(), err, "get user")
		if assert.NotNil(s.T(), actual) {
			assert.Equal(s.T(), "test", actual.Name)
			assert.Equal(s.T(), "pass", actual.Password)
			assert.Equal(s.T(), "zone", actual.Timezone)
		}
	}
}

func (s *TestSuite) Test_CreateUser() {
	user, err := s.repo.CreateUser("create", "user", "test")
	assert.NoError(s.T(), err, "create user")
	if assert.NotNil(s.T(), user) {
		assert.Equal(s.T(), "create", user.Name)
		assert.Equal(s.T(), "user", user.Password)
		assert.Equal(s.T(), "test", user.Timezone)
	}
}

func (s *TestSuite) Test_UpdateUserTimezone() {
	s.repo.CreateUser("update", "user", "test")

	updated, err := s.repo.UpdateUserTimezone("update", "updated")
	assert.NoError(s.T(), err, "update user timezone")
	if assert.NotNil(s.T(), updated) {
		assert.Equal(s.T(), "update", updated.Name)
		assert.Equal(s.T(), "user", updated.Password)
		assert.Equal(s.T(), "updated", updated.Timezone)

		actual, err := s.repo.GetUser("update")
		assert.NoError(s.T(), err, "get user")
		if assert.NotNil(s.T(), actual) {
			assert.Equal(s.T(), "update", actual.Name)
			assert.Equal(s.T(), "user", actual.Password)
			assert.Equal(s.T(), "updated", actual.Timezone)
		}
	}
}

func (s *TestSuite) Test_GetEventOwner() {
	userEvent, err := s.repo.CreateEvent(userName, "t", "d", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1"})
	require.NoError(s.T(), err, "create event")
	require.NotNil(s.T(), userEvent)

	s.repo.CreateUser("other", "other", "other")
	otherEvent, err := s.repo.CreateEvent("other", "t", "d", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1"})
	require.NoError(s.T(), err, "create event")
	require.NotNil(s.T(), otherEvent)

	owner, err := s.repo.GetEventOwner(userEvent.ID)
	assert.NoError(s.T(), err, "get event owner")
	assert.Equal(s.T(), userName, owner)

	owner, err = s.repo.GetEventOwner(otherEvent.ID)
	assert.NoError(s.T(), err, "get event owner")
	assert.Equal(s.T(), "other", owner)
}
