package models

import (
	"errors"
	"example.com/event_booking/db"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func init() {
	db.InitDB("_testing", "memory")
}

func mockDb(code func(sqlmock.Sqlmock)) {
	mockedDb, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	original := db.DB
	db.DB = mockedDb

	defer func() {
		db.DB.Close()
		db.DB = original
	}()

	code(mock)
}

func TestSave(t *testing.T) {
	tests := []struct {
		name     string
		expected func(*Event) (bool, error)
	}{
		{
			name: "It assign newly generated id",
			expected: func(event *Event) (bool, error) {
				error := event.Save()

				return error != nil || event.ID == 0, error
			},
		},
		{
			name: "It return does not update event if query does not prepared",
			expected: func(event *Event) (bool, error) {
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`INSERT INTO events \(name, description, location, date_time, user_id\)`).
						WillReturnError(fmt.Errorf("some error"))
					err = event.Save()
				})

				return err == nil || event.ID != 0, err
			},
		},
		{
			name: "It does not update id when error happens",
			expected: func(event *Event) (bool, error) {
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`INSERT INTO events \(name, description, location, date_time, user_id\)`).
						ExpectExec().
						WithArgs(event.Name, event.Description, event.Location, event.DateTime, event.UserID).
						WillReturnError(fmt.Errorf("some error"))
					err = event.Save()
				})

				return err == nil || event.ID != 0, err
			},
		},
		{
			name: "It return an error if LastInsertedId failed",
			expected: func(event *Event) (bool, error) {
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`INSERT INTO events \(name, description, location, date_time, user_id\)`).
						ExpectExec().
						WithArgs(event.Name, event.Description, event.Location, event.DateTime, event.UserID).
						WillReturnResult(sqlmock.NewErrorResult(errors.New("Some test error")))
					err = event.Save()
				})

				return err == nil || event.ID != 0, err
			},
		},
		{
			name: "When event saved twice",
			expected: func(event *Event) (bool, error) {
				err := event.Save()
				oldId := event.ID
				err = event.Save()

				return err != nil || event.ID == 0 || oldId != event.ID, err
			},
		},
	}

	var event Event

	for _, test := range tests {
		event = Event{
			Name:        "Test name",
			Description: "Test description",
			DateTime:    time.Now(),
			UserID:      4,
		}

		t.Run(test.name, func(t *testing.T) {
			if got, err := test.expected(&event); got {
				t.Errorf("got %v, want %v, error: %v", 0, event.ID, err)
			}
		})
	}
}

func TestFindEvent(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It returs an event from db",
			test: func() (bool, error) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "location", "date_time", "user_id"}).
					AddRow(1, "Test", "Test decription", "USA", time.Now(), 2)
				var event *Event
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(`SELECT \* FROM events WHERE id =`).
						WithArgs(1).
						WillReturnRows(rows)

					event, err = FindEvent(1)
				})

				return !(err == nil && *event != Event{}), err
			},
		},
		{
			name: "It return an error when event does not exist",
			test: func() (bool, error) {
				event, err := FindEvent(0)

				return !(*event == Event{} && err != nil), err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, err := test.test(); got {
				t.Errorf("Test faild with error: %v", err)
			}
		})
	}
}

func TestFindEventByIdUserId(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It return an event when event exist",
			test: func() (bool, error) {
				dateTimeOfEvent, err := time.Parse(time.RFC3339, "2026-05-29T12:19:53+07:00")

				newEvent := Event{
					Name:        "Test",
					Description: "Test desc",
					Location:    "Test",
					DateTime:    dateTimeOfEvent,
					UserID:      3,
				}
				newEvent.Save()

				event, err := FindEventByIdUserId(newEvent.ID, 3)
				return !(newEvent == *event && err == nil), err
			},
		},
		{
			name: "It returns an error when event does not exist",
			test: func() (bool, error) {
				event, err := FindEventByIdUserId(0, 0)

				return !(event == nil || err != nil), err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, err := test.test(); got {
				t.Errorf("Test faild with error: %v", err)
			}
		})
	}
}
