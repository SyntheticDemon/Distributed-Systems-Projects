package internal

import (
	"TicketReservation/pkg/models"
	"fmt"
	"sync"
	"time"
)

type TicketService struct {
	events sync.Map
	lock   sync.Mutex
}

func (ts *TicketService) CreateEvent(name string, date time.Time, totalTickets int) (*models.Event, error) {
	event := &models.Event{
		ID:               generateUUID(),
		Name:             name,
		Date:             date,
		TotalTickets:     totalTickets,
		AvailableTickets: totalTickets,
	}

	ts.events.Store(event.ID, event)
	return event, nil
}

func (ts *TicketService) ListEvents() []*models.Event {
	var events []*models.Event
	ts.events.Range(func(key, value interface{}) bool {
		event := value.(*models.Event)
		events = append(events, event)
		return true
	})
	return events
}

func (ts *TicketService) BookTickets(eventID string, numTickets int) ([]string, error) {

	ts.lock.Lock()

	defer ts.lock.Unlock()

	event, ok := ts.events.Load(eventID)
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	ev := event.(*models.Event)
	if ev.AvailableTickets < numTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	var ticketIDs []string
	for i := 0; i < numTickets; i++ {
		ticketID := generateUUID()
		ticketIDs = append(ticketIDs, ticketID)
	}

	ev.AvailableTickets -= numTickets
	ts.events.Store(eventID, ev)

	return ticketIDs, nil
}
