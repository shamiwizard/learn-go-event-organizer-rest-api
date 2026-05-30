package models

import (
	"errors"
	"example.com/event_booking/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func New(id int64, name string, desc string, date time.Time, userId int64) Event {
	return Event{
		ID:          id,
		Name:        name,
		Description: desc,
		DateTime:    date,
		UserID:      userId,
	}
}

func (event *Event) Save() error {
	var query string
	updateColumn := []any{
		event.Name,
		event.Description,
		event.Location,
		event.DateTime,
		event.UserID,
	}

	if event.ID == 0 {
		query = `INSERT INTO events (name, description, location, date_time, user_id)
		VALUES (?, ?, ?, ?, ?);`
	} else {
		query = `UPDATE events
		SET name = ?,
				description = ?,
				location = ?,
				date_time = ?,
				user_id = ?
		WHERE id = ?
		`
		updateColumn = append(updateColumn, event.ID)
	}

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(updateColumn...)

	if err != nil {
		return err
	}

	generatedId, err := result.LastInsertId()

	event.ID = generatedId

	return err
}

func FindEvent(id int64) (*Event, error) {
	row := db.DB.QueryRow("SELECT * FROM events WHERE id = ?", id)
	event := Event{}

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	return &event, err
}

func FindEventByIdUserId(id, userId int64) (*Event, error) {
	row := db.DB.QueryRow("SELECT * FROM events WHERE id = ? AND user_id = ?", id, userId)
	event := Event{}

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	return &event, err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events;"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, err
}

func (event *Event) Update() error {
	query := `UPDATE events 
						SET name = ?,
								description = ?,
								location = ?,
								date_time = ?,
								user_id = ?
						WHERE id = ?;
						`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		event.Name,
		event.Description,
		event.Location,
		event.DateTime,
		event.UserID,
		event.ID,
	)

	return err
}

func (event *Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?;"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)

	return err
}

func (event *Event) RegisterUser(user *User) error {
	if event.isUserOwner(user) {
		return errors.New("User is an owner of the event")
	}

	if isUserRegisteredToEvent(user, event) {
		return errors.New("User already registered")
	}

	return registerUserToEvent(user, event)
}

func (event *Event) CancelRegistration(user *User) error {
	if event.isUserOwner(user) {
		return errors.New("User is an owner of the event")
	}

	if !isUserRegisteredToEvent(user, event) {
		return errors.New("User is not registered to the event")
	}

	return cancelUserEvetRegistration(user, event)
}

func (event *Event) isUserOwner(user *User) bool {
	if event.UserID == user.ID {
		return true
	}

	return false
}
