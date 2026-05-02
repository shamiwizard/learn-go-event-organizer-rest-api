package models

import (
	"example.com/event_booking/db"
	"example.com/event_booking/utils"
)

type User struct {
	ID int64
	Email string `binding:"required"`
	Password string `binding:"required"`
}


func (user *User) Save() error {
	query := `
	INSERT INTO users (email, password)
	VALUES(?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	user.Password, err = utils.HashPassword(user.Password)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(user.Email, user.Password)
	
	if err != nil {
		return err
	}

	user.ID, err = result.LastInsertId()

	return err	
}
