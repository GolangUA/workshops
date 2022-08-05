package calendar

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
	original := systemLocation
	defer func() {
		systemLocation = original
	}()

	systemLocation = time.UTC

	var err error

	_, _, err = normalizeDateTime("fake date", "11:11", "Europe/Kiev")
	assert.ErrorContains(t, err, `convert datetime="fake date 11:11"`)

	_, _, err = normalizeDateTime("2022-11-22", "fake time", "Europe/Kiev")
	assert.ErrorContains(t, err, `convert datetime="2022-11-22 fake time"`)

	_, _, err = normalizeDateTime("2022-11-22", "01:02", "fake zone")
	assert.ErrorContains(t, err, `invalid timezone="fake zone"`)
}

func Test_normalizeDateTime(t *testing.T) {
	original := systemLocation
	defer func() {
		systemLocation = original
	}()

	systemLocation = time.UTC

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
