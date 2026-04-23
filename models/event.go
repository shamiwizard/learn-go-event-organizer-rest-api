package models

import (
	"time"
)

type Event struct {
	ID int
	Name string `binding:"required"`
	Description string `binding:"required"`
	DateTime time.Time `binding:"required"`
	UserID int
}

func New(id int, name string, desc string, date time.Time, userId int) Event {
	return Event{
		ID: id,
		Name: name,
		Description: desc,
		DateTime: date, 
		UserID: userId,
	}
}

var events = []Event{}

func (event Event) Save() []Event {
	events = append(events, event)

	return events
}

func GetAllEvents() []Event {
	return events
}
