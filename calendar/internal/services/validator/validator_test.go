package validator

import (
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Validate_GetEvents(t *testing.T) {
	service := Service{}

	tests := []struct {
		expected string
		GetEvents
	}{
		{"validation errors: [timezone must be a valid timezone]", GetEvents{"", "invalid", "", "", "", ""}},
		{"validation errors: [date_from must be a valid date]", GetEvents{"", "", "2022/08/10", "", "", ""}},
		{"validation errors: [time_from must be a valid time]", GetEvents{"", "", "", "11:22:33", "", ""}},
		{"validation errors: [date_to must be a valid date]", GetEvents{"", "", "", "", "2022/08/10", ""}},
		{"validation errors: [time_to must be a valid time]", GetEvents{"", "", "", "", "", "11:22:33"}},
		{"validation errors: [both date_from and time_from should be present if timezone is specified]", GetEvents{"", "Europe/Kiev", "2022-01-01", "", "", ""}},
		{"validation errors: [both date_from and time_from should be present if timezone is specified]", GetEvents{"", "Europe/Kiev", "", "12:00", "", ""}},
		{"validation errors: [both date_to and time_to should be present if timezone is specified]", GetEvents{"", "Europe/Kiev", "", "", "2021-01-01", ""}},
		{"validation errors: [both date_to and time_to should be present if timezone is specified]", GetEvents{"", "Europe/Kiev", "", "", "", "18:00"}},
		{"validation errors: [timezone must be a valid timezone; date_from must be a valid date; time_to must be a valid time; both date_from and time_from should be present if timezone is specified; both date_to and time_to should be present if timezone is specified]", GetEvents{"", "invalid", "invalid", "", "", "invalid"}},
		{"validation errors: [timezone must be a valid timezone; time_from must be a valid time; date_to must be a valid date; time_to must be a valid time; both date_from and time_from should be present if timezone is specified]", GetEvents{"", "invalid", "", "invalid", "invalid", "invalid"}},
	}

	for _, test := range tests {
		err := service.Validate(&test.GetEvents)
		assert.Errorf(t, err, "validation should fail for %+v", test)
		assert.EqualErrorf(t, err, test.expected, "error message should be correct")
	}
}

func Test_Validate_CreateEvent(t *testing.T) {
	service := Service{}

	tests := []struct {
		expected string
		CreateEvent
	}{
		{"validation errors: [title must not be blank]", CreateEvent{"", "", "2022-01-01 12:00", "Europe/Kiev", 10, []string{}}},
		{"validation errors: [invalid time value]", CreateEvent{"title", "", "2022-01-01 12:00:00", "Europe/Kiev", 10, []string{}}},
		{"validation errors: [timezone must not be blank and be a valid timezone]", CreateEvent{"title", "", "2022-01-01 12:00", "", 10, []string{}}},
		{"validation errors: [timezone must not be blank and be a valid timezone]", CreateEvent{"title", "", "2022-01-01 12:00", "invalid", 10, []string{}}},
		{"validation errors: [duration must be greater than 0]", CreateEvent{"title", "", "2022-01-01 12:00", "Europe/Kiev", 0, []string{}}},
		{"validation errors: [duration must be greater than 0]", CreateEvent{"title", "", "2022-01-01 12:00", "Europe/Kiev", -1, []string{}}},
	}

	for _, test := range tests {
		err := service.Validate(&test.CreateEvent)
		assert.Errorf(t, err, "validation should fail for %+v", test)
		assert.EqualErrorf(t, err, test.expected, "error message should be correct")
	}
}

func Test_Validate_UpdateEvent(t *testing.T) {
	service := Service{}

	tests := []struct {
		expected string
		UpdateEvent
	}{
		{"validation errors: [id must not be blank]", UpdateEvent{"", "title", "", "2022-01-01 12:00", "Europe/Kiev", 10, []string{}}},
		{"validation errors: [title must not be blank]", UpdateEvent{"EID", "", "", "2022-01-01 12:00", "Europe/Kiev", 10, []string{}}},
		{"validation errors: [invalid time value]", UpdateEvent{"EID", "title", "", "2022-01-01 12:00:00", "Europe/Kiev", 10, []string{}}},
		{"validation errors: [timezone must not be blank and be a valid timezone]", UpdateEvent{"EID", "title", "", "2022-01-01 12:00", "", 10, []string{}}},
		{"validation errors: [timezone must not be blank and be a valid timezone]", UpdateEvent{"EID", "title", "", "2022-01-01 12:00", "invalid", 10, []string{}}},
		{"validation errors: [duration must be greater than 0]", UpdateEvent{"EID", "title", "", "2022-01-01 12:00", "Europe/Kiev", 0, []string{}}},
		{"validation errors: [duration must be greater than 0]", UpdateEvent{"EID", "title", "", "2022-01-01 12:00", "Europe/Kiev", -1, []string{}}},
	}

	for _, test := range tests {
		err := service.Validate(&test.UpdateEvent)
		assert.Errorf(t, err, "validation should fail for %+v", test)
		assert.EqualErrorf(t, err, test.expected, "error message should be correct")
	}
}

func Test_UserTimezone(t *testing.T) {
	service := Service{}

	tests := []struct {
		expected string
		api.UserTimezone
	}{
		{"validation errors: [username must not be blank]", api.UserTimezone{Username: "", Timezone: "Europe/Kiev"}},
		{"validation errors: [timezone must be a valid timezone]", api.UserTimezone{Username: "user", Timezone: "invalid"}},
		{"validation errors: [timezone must not be blank]", api.UserTimezone{Username: "user", Timezone: ""}},
	}

	for _, test := range tests {
		err := service.Validate(&test.UserTimezone)
		assert.Errorf(t, err, "validation should fail for %+v", test)
		assert.EqualErrorf(t, err, test.expected, "error message should be correct")
	}
}
