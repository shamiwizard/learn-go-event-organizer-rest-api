package models

import (
	"example.com/event_booking/db"
)

func registerUserToEvent(user *User, event *Event) error {
	query := `INSERT INTO registrations (event_id, user_id)
						VALUES (?, ?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, user.ID)

	return err
}

func cancelUserEvetRegistration(user *User, event *Event) error {
	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?;`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, user.ID)

	return err
}

func isUserRegisteredToEvent(user *User, event *Event) bool {
	query := "SELECT EXISTS(SELECT 1 FROM registrations WHERE event_id = ? AND user_id = ?);"

	var recordExists bool
	err := db.DB.QueryRow(query, event.ID, user.ID).Scan(&recordExists)

	return err != nil || recordExists
}
