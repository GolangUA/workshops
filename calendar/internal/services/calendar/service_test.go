package calendar

import (
	"errors"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetEvents(t *testing.T) {
	original := dbLocation
	defer func() {
		dbLocation = original
	}()
	dbLocation = time.UTC

	tc1 := []*models.Event{}
	tc2 := []*models.Event{{ID: "TC2"}}
	tc3 := []*models.Event{{ID: "TC3"}}
	tc4 := []*models.Event{{ID: "TC4"}}
	tc5 := []*models.Event{{ID: "TC5"}}
	tc6 := []*models.Event{{ID: "TC6"}}

	assertEvents := func(t *testing.T, expected []*models.Event, actual []*models.Event) {
		assert.Equalf(t, len(expected), len(actual), "expected %d events but actual %d", len(expected), len(actual))
		for i, e := range expected {
			assert.Equalf(t, e.ID, actual[i].ID, "expected event %d but actual %d", i, i)
		}
	}

	type repoKey struct {
		username string
		title    string
		dateFrom string
		timeFrom string
		dateTo   string
		timeTo   string
	}
	type testCase struct {
		username string
		title    string
		dateFrom string
		timeFrom string
		dateTo   string
		timeTo   string
		timezone string
		repoKey  repoKey
		repoResp []*models.Event
	}

	testCases := []testCase{
		{"u1", "t1", "2022-08-02", "", "", "", "Europe/Kiev", repoKey{"u1", "t1", "2022-08-01", "21:00", "", ""}, tc1},
		{"u2", "t2", "", "", "2022-08-02", "22:00", "Europe/Kiev", repoKey{"u2", "t2", "", "", "2022-08-02", "19:00"}, tc2},
		{"u3", "t3", "2022-08-01", "15:00", "2022-08-03", "11:00", "Etc/GMT+2", repoKey{"u3", "t3", "2022-08-01", "17:00", "2022-08-03", "13:00"}, tc3},
		{"u4", "t4", "2022-08-01", "23:00", "", "", "", repoKey{"u4", "t4", "2022-08-02", "01:00", "", ""}, tc4},
		{"u5", "t5", "", "", "2022-08-05", "03:00", "Pacific/Guadalcanal", repoKey{"u5", "t5", "", "", "2022-08-04", "16:00"}, tc5},
		{"u6", "t6", "", "", "2022-08-01", "", "UTC", repoKey{"u6", "t6", "", "", "2022-08-01", "23:59"}, tc6},
	}

	db := make(map[repoKey][]*models.Event)
	for _, tc := range testCases {
		db[tc.repoKey] = tc.repoResp
	}

	mockedRepository := &RepositoryMock{
		GetEventsFunc: func(username string, title string, dateFrom string, timeFrom string, dateTo string, timeTo string) ([]*models.Event, error) {
			if resp, ok := db[repoKey{username, title, dateFrom, timeFrom, dateTo, timeTo}]; ok {
				return resp, nil
			}
			return nil, errors.New("not found")
		},
		GetUserFunc: func(username string) (*models.User, error) {
			switch username {
			case "user_not_found":
				return nil, errors.New("not found")
			case "u4":
				return &models.User{Name: "u4", Timezone: "Etc/GMT+2"}, nil
			default:
				return &models.User{Name: username, Timezone: "Europe/Kiev"}, nil
			}
		},
	}

	service := NewService(mockedRepository)

	for _, tc := range testCases {
		events, err := service.GetEvents(tc.username, tc.title, tc.dateFrom, tc.timeFrom, tc.dateTo, tc.timeTo, tc.timezone)
		assert.NoError(t, err)
		assertEvents(t, tc.repoResp, events)
	}

	_, err := service.GetEvents("user_not_found", "", "", "", "", "", "")
	assert.Error(t, err)
	_, err = service.GetEvents("events_not_found", "", "", "", "", "", "")
	assert.Error(t, err)
}

func Test_GetEvents_UserTimezone(t *testing.T) {
	original := dbLocation
	defer func() {
		dbLocation = original
	}()
	dbLocation = time.UTC

	mockedRepository := &RepositoryMock{
		GetEventsFunc: func(username string, title string, dateFrom string, timeFrom string, dateTo string, timeTo string) ([]*models.Event, error) {
			tf, _ := time.Parse(time.RFC3339, "2022-08-01T09:00:00Z")
			tt, _ := time.Parse(time.RFC3339, "2022-08-01T15:00:00Z")
			return []*models.Event{
				{
					ID:          "EID",
					Title:       "ET",
					Description: "ED",
					TimeFrom:    tf,
					TimeTo:      tt,
					Notes:       []string{},
				}}, nil
		},
		GetUserFunc: func(username string) (*models.User, error) {
			return &models.User{Name: username, Timezone: "Europe/Kiev"}, nil
		},
	}

	service := NewService(mockedRepository)

	events, err := service.GetEvents("user", "", "", "", "", "", "Europe/Kiev")
	assert.NoError(t, err)
	assert.Equalf(t, 1, len(events), "expected 1 event but actual %d", len(events))
	loc, _ := time.LoadLocation("Europe/Kiev")
	tf, _ := time.ParseInLocation(time.RFC3339, "2022-08-01T12:00:00+03:00", loc)
	tt, _ := time.ParseInLocation(time.RFC3339, "2022-08-01T18:00:00+03:00", loc)
	assert.Equal(t, tf, events[0].TimeFrom)
	assert.Equal(t, tt, events[0].TimeTo)
}

func Test_GetEvent(t *testing.T) {
	now := time.Now()
	mockedRepository := &RepositoryMock{
		GetEventFunc: func(id string) (*models.Event, error) {
			return &models.Event{ID: id, Title: "ET", Description: "ED", TimeFrom: now, TimeTo: now, Notes: []string{"note"}}, nil
		},
	}

	service := NewService(mockedRepository)

	event, err := service.GetEvent("EID")
	assert.NoError(t, err)
	assert.Equal(t, models.Event{
		ID:          "EID",
		Title:       "ET",
		Description: "ED",
		TimeFrom:    now,
		TimeTo:      now,
		Notes:       []string{"note"},
	}, *event)
}

func Test_CreateEvent(t *testing.T) {
	mockedRepository := &RepositoryMock{
		CreateEventFunc: func(username string, title string, description string, from time.Time, to time.Time, notes []string) (*models.Event, error) {
			return &models.Event{ID: "EID", Title: title, Description: description, TimeFrom: from, TimeTo: to, Notes: notes}, nil
		},
	}

	event, err := NewService(mockedRepository).
		CreateEvent("username", "title", "description", "2022-02-03 04:56", "UTC", 75*time.Minute, []string{"note"})
	assert.NoError(t, err)
	assert.Equal(t, models.Event{
		ID:          "EID",
		Title:       "title",
		Description: "description",
		TimeFrom:    time.Date(2022, time.February, 3, 4, 56, 0, 0, time.UTC),
		TimeTo:      time.Date(2022, time.February, 3, 6, 11, 0, 0, time.UTC),
		Notes:       []string{"note"},
	}, *event)
}

func Test_UpdateEvent(t *testing.T) {
	mockRepository := &RepositoryMock{
		UpdateEventFunc: func(id string, title string, description string, from time.Time, to time.Time, notes []string) (*models.Event, error) {
			return &models.Event{ID: id, Title: title, Description: description, TimeFrom: from, TimeTo: to, Notes: notes}, nil
		},
		EventExistsFunc: func(id string) (bool, error) {
			if id == "not_exist" {
				return false, nil
			} else if id == "exist_check_error" {
				return false, errors.New("error")
			}
			return true, nil
		},
	}

	service := NewService(mockRepository)

	_, err := service.UpdateEvent("exist_check_error", "title", "description", "2022-02-03 04:56", "UTC", 75*time.Minute, []string{"note"})
	assert.Error(t, err)

	event, err := service.UpdateEvent("not_exist", "title", "description", "2022-02-03 04:56", "UTC", 75*time.Minute, []string{"note"})
	assert.NoErrorf(t, err, "expected no error but actual %s", err)
	assert.Nil(t, event)

	event, err = service.UpdateEvent("exist", "title", "description", "2022-02-03 04:56", "UTC", 75*time.Minute, []string{"note"})
	assert.NoError(t, err)
	assert.Equal(t, models.Event{
		ID:          "exist",
		Title:       "title",
		Description: "description",
		TimeFrom:    time.Date(2022, time.February, 3, 4, 56, 0, 0, time.UTC),
		TimeTo:      time.Date(2022, time.February, 3, 6, 11, 0, 0, time.UTC),
		Notes:       []string{"note"},
	}, *event)
}

func Test_timeFromTo_Errors(t *testing.T) {
	var err error

	_, _, err = timeFromTo("2022-01-23 34:56", "Invalid/Zone", time.Minute*10)
	assert.ErrorContains(t, err, `invalid location "Invalid/Zone"`)

	_, _, err = timeFromTo("2022-02-03 04:50:06", "", time.Minute*20)
	assert.ErrorContains(t, err, `invalid datetime "2022-02-03 04:50:06"`)

	_, _, err = timeFromTo("2022-02-03 04:50:06", "UTC", time.Second*-1)
	assert.Error(t, err, "duration must be greater than 0")

	_, _, err = timeFromTo("2022-03-04 05:06", "Europe/Kiev", time.Second*0)
	assert.Error(t, err, "duration must be greater than 0")
}

func Test_timeFromTo(t *testing.T) {
	type testCase struct {
		expFrom  string
		expTo    string
		datetime string
		timezone string
		duration int
	}

	data := []testCase{
		{"2022-01-02T03:45:00-03:00", "2022-01-02T03:55:00-03:00", "2022-01-02 03:45", "America/Araguaina", 10},
		{"2021-11-22T13:27:00-01:00", "2021-11-22T15:42:00-01:00", "2021-11-22 13:27", "Atlantic/Azores", 135},
		{"2022-05-01T23:57:00+11:00", "2022-05-02T00:04:00+11:00", "2022-05-01 23:57", "Pacific/Guadalcanal", 7},
		{"2022-07-07T07:07:00Z", "2022-07-08T03:41:00Z", "2022-07-07 07:07", "UTC", 1234},
	}

	for _, d := range data {
		actFrom, actTo, err := timeFromTo(d.datetime, d.timezone, time.Minute*time.Duration(d.duration))
		expFrom, _ := time.Parse(time.RFC3339, d.expFrom)
		expTo, _ := time.Parse(time.RFC3339, d.expTo)

		assert.NoErrorf(t, err, `no error expected for datetime="%s" timezone="%s" duration="%d"`, d.datetime, d.timezone, d.duration)
		assert.Equalf(t, expFrom.UTC(), actFrom.UTC(), `datetime_from expected="%s" but actual="%s" for datetime="%s" timezone="%s" duration="%d"`, d.expFrom, d.expTo, d.datetime, d.timezone, d.duration)
		assert.Equalf(t, expTo.UTC(), actTo.UTC(), `datetime_to expected="%s" but actual="%s" for datetime="%s" timezone="%s" duration="%d"`, d.expFrom, d.expTo, d.datetime, d.timezone, d.duration)
	}
}

func Test_normalizeDateTime_Errors(t *testing.T) {
	original := dbLocation
	defer func() {
		dbLocation = original
	}()

	dbLocation = time.UTC

	var err error

	_, _, err = normalizeDateTime("fake date", "11:11", "Europe/Kiev")
	assert.ErrorContains(t, err, `convert datetime="fake date 11:11"`)

	_, _, err = normalizeDateTime("2022-11-22", "fake time", "Europe/Kiev")
	assert.ErrorContains(t, err, `convert datetime="2022-11-22 fake time"`)

	_, _, err = normalizeDateTime("2022-11-22", "01:02", "fake zone")
	assert.ErrorContains(t, err, `invalid timezone="fake zone"`)
}

func Test_normalizeDateTime(t *testing.T) {
	original := dbLocation
	defer func() {
		dbLocation = original
	}()

	dbLocation = time.UTC

	type testCase struct {
		expDate  string
		expTime  string
		date     string
		timev    string
		timezone string
	}

	data := []testCase{
		{"2022-01-01", "08:22", "2022-01-01", "10:22", "Europe/Kiev"},
		{"2022-02-02", "23:55", "2022-02-03", "01:55", "Europe/Kiev"},
		{"2022-06-01", "04:47", "2022-05-31", "19:17", "Pacific/Marquesas"},
		{"2022-05-31", "18:39", "2022-05-31", "09:09", "Pacific/Marquesas"},
		{"2022-02-03", "04:05", "2022-02-03", "04:05", ""},
		{"2022-03-04", "05:06", "2022-03-04", "05:06", "UTC"},
	}

	for _, d := range data {
		actDate, actTime, err := normalizeDateTime(d.date, d.timev, d.timezone)

		assert.NoErrorf(t, err, `no error expected for date="%s" time="%s"`, d.date, d.timev)
		assert.Equalf(t, d.expDate, actDate, `expected date="%s" for input="%s" but actual="%s"`, d.expDate, d.date, actDate)
		assert.Equalf(t, d.expTime, actTime, `expected time="%s" for input="%s" but actual="%s"`, d.expTime, d.timev, actTime)
	}
}
