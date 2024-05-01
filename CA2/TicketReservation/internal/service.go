package internal

import (
	"TicketReservation/pkg/models"
	Distrubted_Systems_git "TicketReservation/proto"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventCacheEntry struct {
	Event  *models.Event
	Expiry time.Time
}

const cacheDuration = 1 * time.Second

func (ts *TicketService) fetchEvent(eventID string) (*models.Event, error) {
	if evt, ok := ts.eventCache.Load(eventID); ok {
		entry := evt.(*EventCacheEntry)
		if time.Now().Before(entry.Expiry) {
			log.Printf("Cache Hit: EventID %s", eventID)
			return entry.Event, nil
		}
		log.Printf("Cache Expired: EventID %s", eventID)
		ts.eventCache.Delete(eventID)
	} else {
		log.Printf("Cache Miss: EventID %s", eventID)
	}

	log.Printf("Fetching from main store: EventID %s", eventID)
	evt, ok := ts.events.Load(eventID)
	if !ok {
		return nil, fmt.Errorf("event not found")
	}
	event := evt.(*models.Event)
	ts.eventCache.Store(eventID, &EventCacheEntry{
		Event:  event,
		Expiry: time.Now().Add(cacheDuration),
	})
	log.Printf("Event re-cached: EventID %s", eventID)

	return event, nil
}

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
			log.Printf("Leaky Bucket: Added token, remaining %d", lb.remaining)
			lb.cond.Broadcast()
		}
		lb.lock.Unlock()
	}
}

func (lb *LeakyBucket) Request(tokens int) bool {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	for lb.remaining < tokens {
		log.Printf("Leaky Bucket: Waiting for tokens, needed %d, available %d", tokens, lb.remaining)
		lb.cond.Wait()
	}
	lb.remaining -= tokens
	log.Printf("Leaky Bucket: Granted %d tokens, remaining %d", tokens, lb.remaining)
	return true
}

type TicketService struct {
	Distrubted_Systems_git.UnimplementedTicketServiceServer
	events      sync.Map
	lock        sync.Mutex
	rateLimiter *LeakyBucket
	eventCache  sync.Map
}

func NewTicketService() *TicketService {
	return &TicketService{
		rateLimiter: NewLeakyBucket(5, time.Second),
	}
}

func (ts *TicketService) CreateEvent(ctx context.Context, req *Distrubted_Systems_git.CreateEventRequest) (*Distrubted_Systems_git.EventResponse, error) {

	ts.lock.Lock()
	defer ts.lock.Unlock()
	if !ts.rateLimiter.Request(1) {
		return nil, status.Error(codes.ResourceExhausted, "Too many requests, please try again later.")
	}
	log.Printf("Creating Event: %v", req)
	event := &models.Event{
		ID:               generateUUID(),
		Name:             req.Name,
		Date:             time.Unix(req.Date, 0),
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
	}

	ts.eventCache.Store(event.ID, &EventCacheEntry{
		Event:  event,
		Expiry: time.Now().Add(cacheDuration),
	})
	ts.events.Store(event.ID, event)
	log.Printf("Event created and cached: ID %s, Name %s", event.ID, event.Name)
	return &Distrubted_Systems_git.EventResponse{
		Id:               event.ID,
		Name:             event.Name,
		Date:             event.Date.Unix(),
		TotalTickets:     event.TotalTickets,
		AvailableTickets: event.AvailableTickets,
	}, nil
}

func (ts *TicketService) BookTickets(ctx context.Context, request *Distrubted_Systems_git.BookRequest) (*Distrubted_Systems_git.BookResponse, error) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	if !ts.rateLimiter.Request(1) {
		return nil, status.Error(codes.ResourceExhausted, "Too many requests, please try again later.")
	}
	event, err := ts.fetchEvent(request.EventId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if event.AvailableTickets < request.NumTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	var ticketIDs []string
	for i := 0; i < int(request.NumTickets); i++ {
		ticketID := generateUUID()
		ticketIDs = append(ticketIDs, ticketID)
	}

	event.AvailableTickets -= request.NumTickets
	log.Printf("Booking Tickets: EventID %s, Tickets Booked %d", request.EventId, request.NumTickets)
	ts.eventCache.Store(request.EventId, &EventCacheEntry{
		Event:  event,
		Expiry: time.Now().Add(cacheDuration),
	})

	ts.events.Store(request.EventId, event)

	return &Distrubted_Systems_git.BookResponse{TicketIds: ticketIDs}, nil
}
func (ts *TicketService) ListEvents(ctx context.Context, req *Distrubted_Systems_git.ListEventsRequest) (*Distrubted_Systems_git.ListEventsResponse, error) {
	var events []*Distrubted_Systems_git.Event
	if !ts.rateLimiter.Request(1) {
		return nil, status.Error(codes.ResourceExhausted, "Too many requests, please try again later.")
	}
	log.Print("Listing all events")
	fetchEventDetails := func(eventID string) (*models.Event, bool) {
		if evt, ok := ts.eventCache.Load(eventID); ok {
			entry := evt.(*EventCacheEntry)
			if time.Now().Before(entry.Expiry) {
				log.Printf("Cache Hit for Event: %s", eventID)
				return entry.Event, true
			}
			log.Printf("Cache Expired for Event: %s, removing from cache", eventID)
			ts.eventCache.Delete(eventID)
		} else {
			log.Printf("Cache Miss for Event: %s", eventID)
		}

		if evt, ok := ts.events.Load(eventID); ok {
			event := evt.(*models.Event)
			ts.eventCache.Store(eventID, &EventCacheEntry{
				Event:  event,
				Expiry: time.Now().Add(cacheDuration),
			})
			log.Printf("Event fetched from main store and re-cached: %s", eventID)
			return event, true
		}
		log.Printf("Event not found in main store: %s", eventID)
		return nil, false
	}

	ts.events.Range(func(key, value interface{}) bool {
		eventID := key.(string)
		event, found := fetchEventDetails(eventID)
		if found {
			events = append(events, &Distrubted_Systems_git.Event{
				Id:               event.ID,
				Name:             event.Name,
				Date:             event.Date.Unix(),
				TotalTickets:     event.TotalTickets,
				AvailableTickets: event.AvailableTickets,
			})
			log.Printf("Event added to response list: %s", event.ID)
		} else {
			log.Printf("Failed to add Event to response list: %s", eventID)
		}
		return true
	})

	log.Printf("Total events listed: %d", len(events))
	return &Distrubted_Systems_git.ListEventsResponse{
		EventDetails: events,
	}, nil
}
