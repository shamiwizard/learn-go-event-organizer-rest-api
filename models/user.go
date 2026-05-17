package models

import (
	"errors"
	"example.com/event_booking/db"
	"example.com/event_booking/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

var CurrentUser *User

func FindByIdEmail(id int64, email string) *User {
	query := "SELECT * FROM users WHERE id = ? AND email = ?;"

	row := db.DB.QueryRow(query, id, email)

	var user User

	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return nil
	}

	return &user
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

func (user *User) ValidateCredentials() error {
	row := db.DB.QueryRow("SELECT id, password FROM users WHERE email = ?;", user.Email)

	var hashedPassword string

	err := row.Scan(&user.ID, &hashedPassword)

	if err != nil {
		return errors.New("Invalid credentials")
	}

	check := utils.CheckPassword(user.Password, hashedPassword)

	if !check {
		return errors.New("Invalid credentials")
	}

	user.Password = hashedPassword

	return err
}
