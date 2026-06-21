package routes

import (
	"encoding/json"
	"errors"
	"example.com/event_booking/db"
	"example.com/event_booking/models"
	"example.com/event_booking/testutils"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"
)

func init() {
	db.InitDB("_testing", "memory")
}

func setupRoutes() *gin.Engine {
	gin.SetMode(gin.TestMode)
	server := gin.Default()

	RegisterRoutes(server)

	return server
}

func createNewEvent() *models.Event {
	event := models.Event{
		Name:        "Test",
		Description: "Test description",
		Location:    "USA",
		DateTime:    time.Now(),
		UserID:      1,
	}

	event.Save()

	return &event
}

func TestGetEvents(t *testing.T) {
	router := setupRoutes()

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "when there are events",
			test: func(t *testing.T) {
				event := createNewEvent()
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/api/v1/events", nil)
				router.ServeHTTP(w, req)

				var requestBody map[string][]models.Event
				err := json.Unmarshal(w.Body.Bytes(), &requestBody)

				if err != nil {
					t.Fatalf("Could not process request body. Error: %v", err)
				}

				createdEventPresent := slices.IndexFunc(requestBody["events"], func(a models.Event) bool {
					return a.ID == event.ID
				})

				if w.Code != http.StatusOK || len(requestBody["events"]) < 1 || createdEventPresent == -1 {
					t.Errorf(
						"want code: %v, got code: %v, want body: %v, got body: %v",
						http.StatusOK, w.Code, map[string][]models.Event{"events": {*event}}, w.Body.String(),
					)
				}
			},
		},
		{
			name: "when there are no events",
			test: func(t *testing.T) {
				_, err := db.DB.Exec("DELETE FROM events;")

				if err != nil {
					t.Fatalf("Could not delete events. Error: %v", err)
				}

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/api/v1/events", nil)
				router.ServeHTTP(w, req)

				stringEvent, _ := json.Marshal(map[string][]models.Event{"events": nil})

				if w.Code != http.StatusOK || w.Body.String() != string(stringEvent) {
					t.Errorf(
						"want code: %v, got code: %v, want body: %v, got body: %v",
						w.Code, http.StatusOK, w.Body.String(), string(stringEvent),
					)
				}
			},
		},
		{
			name: "when query returns an error",
			test: func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/api/v1/events", nil)

				testutils.MockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(`SELECT \* FROM events`).WillReturnError(errors.New("Some db error"))
					router.ServeHTTP(w, req)
				})

				stringEvent, _ := json.Marshal(map[string]any{"message": "Could not find events. Try later."})

				if w.Code != http.StatusInternalServerError || w.Body.String() != string(stringEvent) {
					t.Errorf(
						"want code: %v, got code: %v, want body: %v, got body: %v",
						http.StatusInternalServerError, w.Code, string(stringEvent), w.Body.String(),
					)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func TestGetEvent(t *testing.T) {
	router := setupRoutes()

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "when there are no errors and event is present",
			test: func(t *testing.T) {
				event := createNewEvent()

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/events/%v", event.ID), nil)
				router.ServeHTTP(w, req)

				strEvent, _ := json.Marshal(map[string]models.Event{"event": *event})

				if w.Code != http.StatusOK || w.Body.String() != string(strEvent) {
					t.Errorf(
						"want code: %v, got code: %v, want body: %v, got body: %v",
						http.StatusOK, w.Code, strEvent, w.Body,
					)
				}
			},
		},
		{
			name: "when there are no errors but event missing",
			test: func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/events/%v", 0), nil)
				router.ServeHTTP(w, req)

				if w.Code != http.StatusNotFound || w.Body.String() != "{}" {
					t.Errorf(
						"want code: %v, got code: %v, want body: %v, got body: %v",
						http.StatusNotFound, w.Code, "{}", w.Body,
					)
				}
			},
		},
		{
			name: "when id is invalid",
			test: func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/events/%v", 1.2), nil)
				router.ServeHTTP(w, req)

				var resBody map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resBody)

				if err != nil {
					t.Fatalf("Could not process response body. Error: %v", err)
				}

				if w.Code != http.StatusBadRequest || resBody["message"] != "Could not procces an ID." {
					t.Errorf(
						"want code: %v, got code: %v, want body: %v, got body: %v",
						http.StatusBadRequest, w.Code, "Could not proces an ID.", w.Body,
					)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}
