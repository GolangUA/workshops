//go:build integration

package postgre

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func (s *TestSuite) Test_GetEvents() {
	// Users test data. Shuffled to check ordering in result. And we have two t2 events to check filter by title.
	s.repo.CreateEvent(userName, "t3", "d3", s.dateTime("2022-08-13 13:00"), s.dateTime("2022-08-23 20:00"), []string{"n3"})
	s.repo.CreateEvent(userName, "t2", "d2", s.dateTime("2022-08-12 12:00"), s.dateTime("2022-08-22 19:00"), []string{"n2"})
	s.repo.CreateEvent(userName, "t5", "d5", s.dateTime("2022-08-15 15:00"), s.dateTime("2022-08-25 22:00"), []string{"n5"})
	s.repo.CreateEvent(userName, "t2", "d2", s.dateTime("2022-08-12 12:00"), s.dateTime("2022-08-22 19:00"), []string{"n2"})
	s.repo.CreateEvent(userName, "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1"})
	s.repo.CreateEvent(userName, "t4", "d4", s.dateTime("2022-08-14 14:00"), s.dateTime("2022-08-24 21:00"), []string{"n4"})

	// Non user test data
	s.repo.CreateUser("other", "other", "UTC")
	s.repo.CreateEvent("other", "t3", "d3", s.dateTime("2022-08-13 13:00"), s.dateTime("2022-08-23 20:00"), []string{"n3"})

	// Get with no filters
	events, err := s.repo.GetEvents("", "", "", "", "", "")
	require.NoError(s.T(), err, "get events")
	assert.Empty(s.T(), events, "events")

	// Get with user filter
	events, err = s.repo.GetEvents("user", "", "", "", "", "")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t1", "t2", "t2", "t3", "t4", "t5")

	events, err = s.repo.GetEvents("other", "", "", "", "", "")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t3")

	// Get with date from filter
	events, err = s.repo.GetEvents("user", "", "2022-08-12", "", "", "")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t2", "t2", "t3", "t4", "t5")

	// Get with date from and time to filter
	events, err = s.repo.GetEvents("user", "", "2022-08-12", "", "", "21:30")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t2", "t2", "t3", "t4")

	// Get with date from, time to and title
	events, err = s.repo.GetEvents("user", "t2", "2022-08-12", "", "", "21:30")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t2", "t2")

	// Get with date from, time to and title
	events, err = s.repo.GetEvents("user", "t3", "2022-08-12", "", "", "21:30")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t3")

	// Get with date from and date/time to filter
	events, err = s.repo.GetEvents("user", "", "2022-08-12", "", "2022-08-23", "21:30")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t2", "t2", "t3")

	// Get with date/time from and date/time to filter
	events, err = s.repo.GetEvents("user", "", "2022-08-12", "12:30", "2022-08-23", "21:30")
	assert.NoError(s.T(), err)
	assertEventsSlice(s.T(), events, "t3")

	// Get with no result
	events, err = s.repo.GetEvents("user", "t2", "2022-08-12", "12:30", "2022-08-23", "21:30")
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), events)
}

func (s *TestSuite) Test_GetEvent() {
	event, err := s.repo.CreateEvent(userName, "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1", "n2", "n3"})
	require.NoError(s.T(), err, "create event")

	actual, err := s.repo.GetEvent(event.ID)
	require.NoError(s.T(), err, "get event")
	assert.Equal(s.T(), "t1", actual.Title)
	assert.Equal(s.T(), "d1", actual.Description)
	assert.Equal(s.T(), s.dateTime("2022-08-11 11:00").UTC(), actual.TimeFrom.UTC())
	assert.Equal(s.T(), s.dateTime("2022-08-21 18:00").UTC(), actual.TimeTo.UTC())
	assert.Equal(s.T(), []string{"n1", "n2", "n3"}, actual.Notes)
}

func (s *TestSuite) Test_EventExists() {
	event, err := s.repo.CreateEvent(userName, "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1", "n2", "n3"})
	require.NoError(s.T(), err, "create event")

	exists, err := s.repo.EventExists(event.ID)
	require.NoError(s.T(), err, "event exists")
	assert.True(s.T(), exists)

	exists, err = s.repo.EventExists("other")
	require.NoError(s.T(), err, "event exists")
	assert.False(s.T(), exists)

	event, err = s.repo.CreateEvent("non_existing", "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1", "n2", "n3"})
	assert.Nil(s.T(), event)
	assert.Error(s.T(), err)
}

func (s *TestSuite) Test_CreateEvent() {
	event, err := s.repo.CreateEvent(userName, "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1"})
	require.NoError(s.T(), err, "create event")
	assert.NotEmpty(s.T(), event.ID)
	assert.Equal(s.T(), "t1", event.Title)
	assert.Equal(s.T(), "d1", event.Description)
	assert.Equal(s.T(), s.dateTime("2022-08-11 11:00").UTC(), event.TimeFrom.UTC())
	assert.Equal(s.T(), s.dateTime("2022-08-21 18:00").UTC(), event.TimeTo.UTC())
	assert.Equal(s.T(), []string{"n1"}, event.Notes)
}

func (s *TestSuite) Test_UpdateEvent() {
	original, err := s.repo.CreateEvent(userName, "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1"})
	require.NoError(s.T(), err, "create event")

	updated, err := s.repo.UpdateEvent(original.ID, "t2", "d2", s.dateTime("2022-08-12 11:00"), s.dateTime("2022-08-22 18:00"), []string{"n2", "n3"})
	require.NoError(s.T(), err, "update event")
	assert.Equal(s.T(), original.ID, updated.ID)
	assert.Equal(s.T(), "t2", updated.Title)
	assert.Equal(s.T(), "d2", updated.Description)
	assert.Equal(s.T(), s.dateTime("2022-08-12 11:00").UTC(), updated.TimeFrom.UTC())
	assert.Equal(s.T(), s.dateTime("2022-08-22 18:00").UTC(), updated.TimeTo.UTC())
	assert.Equal(s.T(), []string{"n2", "n3"}, updated.Notes)

	actual, err := s.repo.GetEvent(original.ID)
	require.NoError(s.T(), err, "get event")
	assert.Equal(s.T(), original.ID, actual.ID)
	assert.Equal(s.T(), "t2", actual.Title)
	assert.Equal(s.T(), "d2", actual.Description)
	assert.Equal(s.T(), s.dateTime("2022-08-12 11:00").UTC(), actual.TimeFrom.UTC())
	assert.Equal(s.T(), s.dateTime("2022-08-22 18:00").UTC(), actual.TimeTo.UTC())
	assert.Equal(s.T(), []string{"n2", "n3"}, actual.Notes)
}

func (s *TestSuite) Test_DeleteEvent() {
	event, err := s.repo.CreateEvent(userName, "t1", "d1", s.dateTime("2022-08-11 11:00"), s.dateTime("2022-08-21 18:00"), []string{"n1"})
	require.NoError(s.T(), err, "create event")

	actual, err := s.repo.GetEvent(event.ID)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), actual)

	ok, err := s.repo.DeleteEvent(event.ID)
	require.NoError(s.T(), err, "delete event")
	assert.True(s.T(), ok)

	actual, err = s.repo.GetEvent(event.ID)
	assert.NoError(s.T(), err, "get event")
	assert.Nil(s.T(), actual)

	ok, err = s.repo.DeleteEvent(event.ID)
	assert.NoError(s.T(), err, "delete event")
	assert.False(s.T(), ok)
}

func assertEventsSlice(t *testing.T, events []*models.Event, titles ...string) {
	if assert.Len(t, events, len(titles)) {
		for i, title := range titles {
			assert.Equalf(t, title, events[i].Title, "title %d to be %s but is %s", i, title, events[i].Title)
		}
	}
}
