package models

import (
	"time"
	"example.com/event_booking/db"
)

type Event struct {
	ID int64
	Name string `binding:"required"`
	Description string `binding:"required"`
	Location string `binding:"required"`
	DateTime time.Time `binding:"required"`
	UserID int
}

func New(id int64, name string, desc string, date time.Time, userId int) Event {
	return Event{
		ID: id,
		Name: name,
		Description: desc,
		DateTime: date, 
		UserID: userId,
	}
}

var events = []Event{}

func (event *Event) Save() error {
	query := `INSERT INTO events (name, description, location, date_time, user_id)
						VALUES (?, ?, ?, ?, ?);`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.UserID)
	 
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

	if err != nil { 
		return nil, err
	}

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
		return  err
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

func (event * Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?;"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)

	return err
}
