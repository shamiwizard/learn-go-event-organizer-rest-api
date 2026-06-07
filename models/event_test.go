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

func newEvent(index int) Event {
	return Event{
		Name:        fmt.Sprint("Test name", index),
		Description: "Test description",
		Location:    "USA",
		DateTime:    time.Now(),
		UserID:      1,
	}
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
			name: "It return an error when LastInsertedId failed",
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
		event = newEvent(1)

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

func TestGetAllEvents(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "it returns all events in DB",
			test: func() (bool, error) {
				for index := range 2 {
					event := newEvent(index)
					event.Save()
				}

				events, err := GetAllEvents()

				return !(len(events) >= 2 && err == nil), err
			},
		},
		{
			name: "It return empty slice when there is not events",
			test: func() (bool, error) {
				_, err := db.DB.Exec("DELETE FROM events;")

				if err != nil {
					return true, err
				}

				events, err := GetAllEvents()

				return !(len(events) == 0 && err == nil), err
			},
		},
		{
			name: "It return nil and error when query failed",
			test: func() (bool, error) {
				var events []Event
				var err error
				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(`SELECT \* FROM events;`).
						WillReturnError(fmt.Errorf("some error"))

					events, err = GetAllEvents()
				})

				return !(len(events) == 0 && err != nil), err
			},
		},
		{
			name: "It return nil and error when one of the row falied",
			test: func() (bool, error) {
				var events []Event
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "name", "description", "location", "date_time", "user_id"}).
						AddRow(1, "Test", "Test decription", "USA", time.Now(), 2).
						AddRow(2, "Test", "Test decription", "USA", "error date", 2)

					mock.ExpectQuery(`SELECT \* FROM events;`).WillReturnRows(rows)

					events, err = GetAllEvents()
				})

				return !(len(events) == 0 && err != nil), err
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

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It stores event in db if there is no error",
			test: func() (bool, error) {
				newEvent := newEvent(1)

				newEvent.Create()

				event, err := FindEvent(newEvent.ID)

				return !(event != nil && err == nil), err
			},
		},
		{
			name: "It returns an error if event has an id",
			test: func() (bool, error) {
				event := newEvent(1)
				event.ID = 1

				err := event.Create()

				return err.Error() != "Event already stored in DB", err
			},
		},
		{
			name: "It return an error when query preparetion returns an error",
			test: func() (bool, error) {
				var event Event
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`INSERT INTO events \(name, description, location, date_time, user_id\)`).
						WillReturnError(fmt.Errorf("some error"))

					event = newEvent(1)
					err = event.Create()
				})

				return !(event.ID == 0 && err.Error() == "some error"), err
			},
		},
		{
			name: "It return an error when query failed to execute",
			test: func() (bool, error) {
				var event Event
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					event = newEvent(1)

					mock.ExpectPrepare(`INSERT INTO events \(name, description, location, date_time, user_id\)`).
						ExpectExec().
						WithArgs(event.Name, event.Description, event.Location, event.DateTime, event.UserID).
						WillReturnError(fmt.Errorf("some error"))

					err = event.Create()
				})

				return !(event.ID == 0 && err.Error() == "some error"), err
			},
		},
		{
			name: "It return an error when insert result does not have last insert id",
			test: func() (bool, error) {
				var event Event
				var err error

				mockDb(func(mock sqlmock.Sqlmock) {
					event = newEvent(1)

					mock.ExpectPrepare(`INSERT INTO events \(name, description, location, date_time, user_id\)`).
						ExpectExec().
						WithArgs(event.Name, event.Description, event.Location, event.DateTime, event.UserID).
						WillReturnResult(sqlmock.NewErrorResult(errors.New("some error")))

					err = event.Create()
				})

				return !(event.ID == 0 && err.Error() == "some error"), err
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

func TestUpdate(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It updates event",
			test: func() (bool, error) {
				newEvent := newEvent(1)
				err := newEvent.Create()

				newName := "Update name"
				newEvent.Name = newName

				err = newEvent.Update()

				event, err := FindEvent(newEvent.ID)

				return !(event.Name == newName && err == nil), err
			},
		},
		{
			name: "It returns an error when event have incorrect id",
			test: func() (bool, error) {
				event := newEvent(1)

				err := event.Update()

				return !(err.Error() == "The id is incorrect"), err
			},
		},
		{
			name: "It returns an error when query could not be prepared",
			test: func() (bool, error) {
				var err error
				event := newEvent(1)
				event.ID = 1

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`UPDATE events SET (.*) WHERE id = ?`).
						WillReturnError(fmt.Errorf("some error"))
					err = event.Update()
				})

				return !(err.Error() == "some error"), err
			},
		},
		{
			name: "It returns an error when query could not be executed",
			test: func() (bool, error) {
				var err error
				event := newEvent(1)
				event.ID = 1

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`UPDATE events SET (.*) WHERE id = ?`).
						ExpectExec().
						WithArgs(event.Name, event.Description, event.Location, event.DateTime, event.UserID, event.ID).
						WillReturnError(fmt.Errorf("some error"))

					err = event.Update()
				})

				return !(err.Error() == "some error"), err
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

func TestDelete(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It delete event when there is no error",
			test: func() (bool, error) {
				deletedEvent := newEvent(1)
				deletedEvent.Create()
				eventId := deletedEvent.ID
				deleteErr := deletedEvent.Delete()

				event, findErr := FindEvent(eventId)
				return !(*event == Event{} && deleteErr == nil && findErr.Error() == "sql: no rows in result set"), deleteErr
			},
		},
		{
			name: "It returns an error when query could not be prepared",
			test: func() (bool, error) {
				var err error
				event := newEvent(1)
				event.ID = 1

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`DELETE FROM events WHERE id =`).
						WillReturnError(fmt.Errorf("some error"))
					err = event.Delete()
				})

				return !(err.Error() == "some error"), err
			},
		},
		{
			name: "It returns an error when query could not be executed",
			test: func() (bool, error) {
				var err error
				event := newEvent(1)
				event.ID = 1

				mockDb(func(mock sqlmock.Sqlmock) {
					mock.ExpectPrepare(`DELETE FROM events WHERE id =`).
						ExpectExec().
						WithArgs(event.ID).
						WillReturnError(fmt.Errorf("some error"))

					err = event.Delete()
				})

				return !(err.Error() == "some error"), err
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

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It creates new registration record",
			test: func() (bool, error) {
				event := newEvent(1)
				event.UserID = 0
				user := User{
					Email:    "test@mail.com",
					Password: "test12345",
				}

				err := event.Save()
				err = user.Save()
				err = event.RegisterUser(&user)

				return !(err == nil), err
			},
		},
		{
			name: "It returns an error when user already registered",
			test: func() (bool, error) {
				event := newEvent(1)
				user := User{
					Email:    "test@mail.com",
					Password: "test12345",
				}

				err := event.Save()
				err = user.Save()

				err = event.RegisterUser(&user)
				err = event.RegisterUser(&user)

				return !(err.Error() == "User already registered"), err
			},
		},
		{
			name: "It returns an error when user is an owner of the event",
			test: func() (bool, error) {
				event := newEvent(1)
				user := User{
					Email:    "test@mail.com",
					Password: "test12345",
				}

				err := user.Save()
				event.UserID = user.ID

				err = event.Save()

				err = event.RegisterUser(&user)

				return !(err.Error() == "User is an owner of the event"), err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, err := test.test(); got {
				t.Errorf("Test faild with error: %v", err)
			}

			db.DB.Exec("DELETE FROM users;")
		})
	}
}

func TestCancelRegistration(t *testing.T) {
	tests := []struct {
		name string
		test func() (bool, error)
	}{
		{
			name: "It creates cancel registration",
			test: func() (bool, error) {
				event := newEvent(1)
				event.UserID = 0
				user := User{
					Email:    "test@mail.com",
					Password: "test12345",
				}

				err := event.Save()
				err = user.Save()
				err = event.RegisterUser(&user)
				err = event.CancelRegistration(&user)

				return !(err == nil), err
			},
		},
		{
			name: "It returns an error when user is not registered yet",
			test: func() (bool, error) {
				event := newEvent(1)
				user := User{
					Email:    "test@mail.com",
					Password: "test12345",
				}

				err := event.Save()
				err = user.Save()

				err = event.CancelRegistration(&user)

				return !(err.Error() == "User is not registered to the event"), err
			},
		},
		{
			name: "It returns an error when user is an owner of the event",
			test: func() (bool, error) {
				event := newEvent(1)
				user := User{
					Email:    "test@mail.com",
					Password: "test12345",
				}

				err := user.Save()
				event.UserID = user.ID

				err = event.Save()

				err = event.CancelRegistration(&user)

				return !(err.Error() == "User is an owner of the event"), err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, err := test.test(); got {
				t.Errorf("Test faild with error: %v", err)
			}

			db.DB.Exec("DELETE FROM users;")
		})
	}
}

func TestIsUserOwner(t *testing.T) {
	tests := []struct {
		name string
		want bool
		got  func() bool
	}{
		{
			name: "It return true when event.UserID is the same as user.ID",
			want: true,
			got: func() bool {
				event := newEvent(1)
				return event.isUserOwner(&User{ID: 1})
			},
		},
		{
			name: "It return false when event.UserID is different from user.ID",
			want: false,
			got: func() bool {
				event := newEvent(1)
				return event.isUserOwner(&User{ID: 5})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := test.want
			got := test.got()

			if want != got {
				t.Errorf("Failed test. want: %v; got: %v", want, got)
			}
		})
	}
}
