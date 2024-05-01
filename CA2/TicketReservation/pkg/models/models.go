package models

import "time"

type Event struct {
	ID               string
	Name             string
	Date             time.Time
	TotalTickets     int32
	AvailableTickets int32
}

type Ticket struct {
	ID      string
	EventID string
}
