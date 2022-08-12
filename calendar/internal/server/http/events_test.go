package http

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_GetEvents(t *testing.T) {
	validator := &ValidatorMock{
		ValidateFunc: func(v interface{}) error {
			return fmt.Errorf("validation error")
		},
	}
	service := &ServiceMock{
		GetEventsFunc: func(username string, title string, dateFrom string, timeFrom string, dateTo string, timeTo string, timezone string) ([]*models.Event, error) {
			switch title {
			case "title":
				if timezone == "Europe/Kiev" {
					return []*models.Event{{
						ID:          "tc",
						Title:       "title",
						Description: "desc",
						TimeFrom:    time.Date(2022, time.March, 7, 5, 24, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
						TimeTo:      time.Date(2022, time.March, 7, 12, 10, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
						Notes:       []string{"note1", "note2"},
					}}, nil
				} else {
					return []*models.Event{{
						ID:          "tc",
						Title:       "title",
						Description: "desc",
						TimeFrom:    time.Date(2022, time.October, 3, 2, 11, 0, 0, time.UTC),
						TimeTo:      time.Date(2022, time.October, 3, 2, 55, 0, 0, time.UTC),
						Notes:       []string{"note1", "note2"},
					}}, nil
				}
			case "error":
				return nil, fmt.Errorf("service error")
			default:
				return nil, nil
			}
		},
	}
	server := NewServer(service, validator, &AuthMock{})

	w, c := authenticatedContext()

	// Validation check
	server.GetEvents(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"validation error"}`, w.Body.String())

	validator.ValidateFunc = func(v interface{}) error { return nil }

	// GET events error
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/api/events?title=error", nil)
	server.GetEvents(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"service error"}`, w.Body.String())

	// GET events happy path
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/api/events?title=title", nil)
	server.GetEvents(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `[{"id":"tc","title":"title","description":"desc","time":"2022-03-07 05:24","timezone":"Europe/Kiev","duration":406,"notes":["note1","note2"]}]`, w.Body.String())

	// GET events specify timzeon
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/api/events?title=title&timezone=UTC", nil)
	server.GetEvents(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `[{"id":"tc","title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}]`, w.Body.String())
}

func Test_GetEvent(t *testing.T) {
	event := &models.Event{
		ID:          "owns",
		Title:       "title",
		Description: "desc",
		TimeFrom:    time.Date(2022, time.March, 7, 5, 24, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
		TimeTo:      time.Date(2022, time.March, 7, 12, 10, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
		Notes:       []string{"note1", "note2", "note3"},
	}
	service := &ServiceMock{
		GetEventFunc: func(id string) (*models.Event, error) {
			if id == "owns" {
				return event, nil
			}
			return nil, nil
		},
		GetEventOwnerFunc: func(id string) (string, error) {
			switch id {
			case "owns":
				return contextUser, nil
			case "not_owns":
				return "other", nil
			case "not_found":
				return "", nil
			case "own_check_error":
				return "", fmt.Errorf("get owner error")
			default:
				panic("unexpected id=" + id)
			}
		},
	}
	server := NewServer(service, &ValidatorMock{}, &AuthMock{})

	w, c := authenticatedContext()

	// GET event -> owner check error
	c.Request = httptest.NewRequest("GET", "/events/own_check_error", nil)
	c.Params = []gin.Param{{Key: "id", Value: "own_check_error"}}
	server.GetEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"get owner error"}`, w.Body.String())

	// GET event -> owner not found
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/event/not_found", nil)
	c.Params = []gin.Param{{Key: "id", Value: "not_found"}}
	server.GetEvent(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"message":"event with ID=\"not_found\" not found"}`, w.Body.String())

	// GET event -> does not own
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/event/not_owns", nil)
	c.Params = []gin.Param{{Key: "id", Value: "not_owns"}}
	server.GetEvent(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, `{"message":"event with ID=\"not_owns\" access denied"}`, w.Body.String())

	// GET event -> happy path
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/event/owns", nil)
	c.Params = []gin.Param{{Key: "id", Value: "owns"}}
	server.GetEvent(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":"owns","title":"title","description":"desc","time":"2022-03-07 05:24","timezone":"Europe/Kiev","duration":406,"notes":["note1","note2","note3"]}`, w.Body.String())

	// GET event -> get event error
	service.GetEventFunc = func(id string) (*models.Event, error) { return nil, fmt.Errorf("get event error") }
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("GET", "/event/owns", nil)
	c.Params = []gin.Param{{Key: "id", Value: "owns"}}
	server.GetEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"get event error"}`, w.Body.String())
}

func Test_PostEvent(t *testing.T) {
	validator := &ValidatorMock{
		ValidateFunc: func(v interface{}) error { return fmt.Errorf("validation error") },
	}
	service := &ServiceMock{
		CreateEventFunc: func(username string, title string, description string, timeVal string, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
			if title == "error" {
				return nil, fmt.Errorf("create event error")
			}

			return &models.Event{
				ID:          "id",
				Title:       title,
				Description: description,
				TimeFrom:    time.Date(2022, time.March, 7, 5, 24, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
				TimeTo:      time.Date(2022, time.March, 7, 12, 10, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
				Notes:       notes,
			}, nil
		},
	}
	server := NewServer(service, validator, &AuthMock{})

	// POST event -> validation error
	w, c := authenticatedContext()
	c.Request = httptest.NewRequest("POST", "/events", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	server.PostEvent(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"validation error"}`, w.Body.String())

	validator.ValidateFunc = func(v interface{}) error { return nil }

	// POST event -> create event error
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("POST", "/events", strings.NewReader(`{"title":"error","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	server.PostEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"create event error"}`, w.Body.String())

	// POST event -> happy path
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("POST", "/events", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1"]}`))
	server.PostEvent(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"id":"id","title":"title","description":"desc","time":"2022-03-07 05:24","timezone":"Europe/Kiev","duration":406,"notes":["note1"]}`, w.Body.String())
}

func Test_PutEvent(t *testing.T) {
	validator := &ValidatorMock{
		ValidateFunc: func(v interface{}) error { return fmt.Errorf("validation error") },
	}
	service := &ServiceMock{
		UpdateEventFunc: func(id string, title string, description string, timeVal string, timezone string, duration time.Duration, notes []string) (*models.Event, error) {
			if id == "update_error" {
				return nil, fmt.Errorf("update event error")
			}

			return &models.Event{
				ID:          "id",
				Title:       title,
				Description: description,
				TimeFrom:    time.Date(2022, time.March, 7, 5, 24, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
				TimeTo:      time.Date(2022, time.March, 7, 12, 10, 0, 0, time.FixedZone("Europe/Kiev", 3*60*60)),
				Notes:       notes,
			}, nil
		},
		GetEventOwnerFunc: func(id string) (string, error) {
			switch id {
			case "owns", "update_error":
				return contextUser, nil
			case "not_owns":
				return "other", nil
			case "not_found":
				return "", nil
			case "own_check_error":
				return "", fmt.Errorf("get owner error")
			default:
				panic("unexpected id=" + id)
			}
		},
	}
	server := NewServer(service, validator, &AuthMock{})

	// PUT event -> validation error
	w, c := authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/id", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "id"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"message":"validation error"}`, w.Body.String())

	validator.ValidateFunc = func(v interface{}) error { return nil }

	// PUT event -> update event error
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/id", strings.NewReader(`{"title":"error","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "update_error"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"update event error"}`, w.Body.String())

	// PUT event -> owner check error
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/id", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "own_check_error"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"get owner error"}`, w.Body.String())

	// PUT event -> does not own event
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/id", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "not_owns"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, `{"message":"event with ID=\"not_owns\" access denied"}`, w.Body.String())

	// PUT event -> owner not found
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/not_found", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "not_found"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"message":"event with ID=\"not_found\" not found"}`, w.Body.String())

	// PUT event -> update error
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/id", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "update_error"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"update event error"}`, w.Body.String())

	// PUT event -> happy path
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("PUT", "/event/id", strings.NewReader(`{"title":"title","description":"desc","time":"2022-10-03 02:11","timezone":"UTC","duration":44,"notes":["note1","note2"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "owns"}}
	server.PutEvent(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":"id","title":"title","description":"desc","time":"2022-03-07 05:24","timezone":"Europe/Kiev","duration":406,"notes":["note1","note2"]}`, w.Body.String())
}

func Test_DeleteEvent(t *testing.T) {
	service := &ServiceMock{
		DeleteEventFunc: func(id string) (bool, error) {
			if id == "owns" {
				return true, nil
			}
			return false, nil
		},
		GetEventOwnerFunc: func(id string) (string, error) {
			switch id {
			case "owns":
				return contextUser, nil
			case "not_owns":
				return "other", nil
			case "not_found":
				return "", nil
			case "own_check_error":
				return "", fmt.Errorf("get owner error")
			default:
				panic("unexpected id=" + id)
			}
		},
	}
	server := NewServer(service, &ValidatorMock{}, &AuthMock{})

	w, c := authenticatedContext()

	// DELETE event -> owner check error
	c.Request = httptest.NewRequest("DELETE", "/events/own_check_error", nil)
	c.Params = []gin.Param{{Key: "id", Value: "own_check_error"}}
	server.DeleteEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"get owner error"}`, w.Body.String())

	// DELETE event -> owner not found
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("DELETE", "/event/not_found", nil)
	c.Params = []gin.Param{{Key: "id", Value: "not_found"}}
	server.DeleteEvent(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"message":"event with ID=\"not_found\" not found"}`, w.Body.String())

	// DELETE event -> does not own
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("DELETE", "/event/not_owns", nil)
	c.Params = []gin.Param{{Key: "id", Value: "not_owns"}}
	server.DeleteEvent(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, `{"message":"event with ID=\"not_owns\" access denied"}`, w.Body.String())

	// DELETE event -> happy path
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("DELETE", "/event/owns", nil)
	c.Params = []gin.Param{{Key: "id", Value: "owns"}}
	server.DeleteEvent(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// DELETE event -> delete event error
	service.DeleteEventFunc = func(id string) (bool, error) { return false, fmt.Errorf("delete event error") }
	w, c = authenticatedContext()
	c.Request = httptest.NewRequest("DELETE", "/event/owns", nil)
	c.Params = []gin.Param{{Key: "id", Value: "owns"}}
	server.DeleteEvent(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `{"message":"delete event error"}`, w.Body.String())
}
