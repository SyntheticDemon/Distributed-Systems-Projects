package internal

import (
	"TicketReservation/pkg/models"
	Distrubted_Systems_git "TicketReservation/proto"
	"context"
	"fmt"
	"sync"
	"time"
)

type TicketService struct {
	Distrubted_Systems_git.UnimplementedTicketServiceServer
	events sync.Map
	lock   sync.Mutex
}

func (ts *TicketService) BookTickets(ctx context.Context, request *Distrubted_Systems_git.BookRequest) (*Distrubted_Systems_git.BookResponse, error) {

	ts.lock.Lock()

	defer ts.lock.Unlock()

	event, ok := ts.events.Load(request.EventId)
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	ev := event.(*models.Event)
	if ev.AvailableTickets < request.NumTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	var ticketIDs []string
	for i := 0; i < int(request.NumTickets); i++ {
		ticketID := generateUUID()
		ticketIDs = append(ticketIDs, ticketID)
	}

	ev.AvailableTickets -= request.NumTickets
	ts.events.Store(request.EventId, ev)

	return &Distrubted_Systems_git.BookResponse{TicketIds: ticketIDs}, nil
}

func (ts *TicketService) CreateEvent(ctx context.Context, req *Distrubted_Systems_git.CreateEventRequest) (*Distrubted_Systems_git.EventResponse, error) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	event := &models.Event{
		ID:               generateUUID(),
		Name:             req.Name,
		Date:             time.Unix(req.Date, 0),
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
	}

	ts.events.Store(event.ID, event)
	return &Distrubted_Systems_git.EventResponse{
		Id:               event.ID,
		Name:             event.Name,
		Date:             event.Date.Unix(),
		TotalTickets:     int32(event.TotalTickets),
		AvailableTickets: int32(event.AvailableTickets),
	}, nil
}

func (ts *TicketService) ListEvents(ctx context.Context, req *Distrubted_Systems_git.ListEventsRequest) (*Distrubted_Systems_git.ListEventsResponse, error) {
	var events []*Distrubted_Systems_git.Event
	ts.events.Range(func(key, value interface{}) bool {
		event := value.(*models.Event)
		events = append(events, &Distrubted_Systems_git.Event{
			Id:               event.ID,
			Name:             event.Name,
			Date:             event.Date.Unix(),
			TotalTickets:     int32(event.TotalTickets),
			AvailableTickets: int32(event.AvailableTickets),
		})
		return true
	})
	return &Distrubted_Systems_git.ListEventsResponse{
		EventDetails: events,
	}, nil
}
