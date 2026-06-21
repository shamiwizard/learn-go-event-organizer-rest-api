package testutils

import (
	"example.com/event_booking/db"
	"github.com/DATA-DOG/go-sqlmock"
)

func MockDb(code func(sqlmock.Sqlmock)) {
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
