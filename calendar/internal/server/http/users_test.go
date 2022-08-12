package http

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_UpdateUserTimezone(t *testing.T) {
	validator := &ValidatorMock{
		ValidateFunc: func(v interface{}) error {
			return fmt.Errorf("validation error")
		},
	}
	service := &ServiceMock{
		UpdateUserTimezoneFunc: func(username string, timezone string) (*models.User, error) {
			if timezone == "UTC" {
				return nil, fmt.Errorf("error")
			}
			return &models.User{
				Name:     username,
				Timezone: timezone,
			}, nil
		},
	}

	server := NewServer(service, validator, &AuthMock{})

	// PUT user -> validation error
	w, c := authenticatedContext()
	c.Request = httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(`{"username":"user","timezone":"Europe/Kyiv"}`))
	server.UpdateUserTimezone(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"validation error"}`, w.Body.String())

	validator.ValidateFunc = func(v interface{}) error { return nil }

	// PUT user -> bad payload
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(`non json payload`))
	server.UpdateUserTimezone(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"failed to parse request body"}`, w.Body.String())

	// PUT user -> forbidden
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(`{"username":"other","timezone":"Europe/Kyiv"}`))
	server.UpdateUserTimezone(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, `{"message":"You are not allowed to update this user access denied"}`, w.Body.String())

	// PUT user -> error
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(`{"username":"user","timezone":"UTC"}`))
	server.UpdateUserTimezone(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"error"}`, w.Body.String())

	// PUT user -> success
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(`{"username":"user","timezone":"Europe/Kyiv"}`))
	server.UpdateUserTimezone(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"username":"user","timezone":"Europe/Kyiv"}`, w.Body.String())
}
