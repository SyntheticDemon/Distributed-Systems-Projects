package models

import "time"

type Event struct {
	ID               string
	Name             string
	Date             time.Time
	TotalTickets     int
	AvailableTickets int
}

type Ticket struct {
	ID      string
	EventID string
}
