package internal

import (
	"TicketReservation/pkg/models"
	Distrubted_Systems_git "TicketReservation/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity  int
	remaining int
	rate      time.Duration
	lock      sync.Mutex
	cond      *sync.Cond
}

func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity:  capacity,
		remaining: capacity,
		rate:      rate,
	}
	lb.cond = sync.NewCond(&lb.lock)
	go lb.leak()
	return lb
}

func (lb *LeakyBucket) leak() {
	for {
		time.Sleep(lb.rate)
		lb.lock.Lock()
		if lb.remaining < lb.capacity {
			lb.remaining++
			lb.cond.Broadcast()
		}
		lb.lock.Unlock()
	}
}

func (lb *LeakyBucket) Request(tokens int) bool {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	for lb.remaining < tokens {
		lb.cond.Wait()
	}
	lb.remaining -= tokens
	return true
}

type TicketService struct {
	Distrubted_Systems_git.UnimplementedTicketServiceServer
	events      sync.Map
	lock        sync.Mutex
	rateLimiter *LeakyBucket
}

func NewTicketService() *TicketService {
	return &TicketService{
		rateLimiter: NewLeakyBucket(5, time.Second), // Allow max 5 requests per second
	}
}

func (ts *TicketService) BookTickets(ctx context.Context, request *Distrubted_Systems_git.BookRequest) (*Distrubted_Systems_git.BookResponse, error) {
	if !ts.rateLimiter.Request(1) {
		return nil, status.Error(codes.ResourceExhausted, "Too many requests, please try again later.")
	}

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
