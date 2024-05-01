package main

import (
	ticket "TicketReservation/proto" // import the generated gRPC client package
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8998", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := ticket.NewTicketServiceClient(conn)

	// Call CreateEvent
	eventId := createEvent(client)

	//// Call ListEvents
	listEvents(client)
	//
	//// Call BookTickets
	bookTickets(client, eventId, 2)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Call BookTickets concurrently for different clients
			bookTickets(client, eventId, 2)
		}(i)
	}

	wg.Wait()
}

func createEvent(client ticket.TicketServiceClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.CreateEvent(ctx, &ticket.CreateEventRequest{
		Name:         "Concert",
		Date:         time.Now().Unix(),
		TotalTickets: 100,
	})
	if err != nil {
		log.Fatalf("Could not create event: %v", err)
	}
	fmt.Printf("Created Event: %v\n", r)
	return r.Id
}

func listEvents(client ticket.TicketServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.ListEvents(ctx, &ticket.ListEventsRequest{})
	if err != nil {
		log.Fatalf("Could not list events: %v", err)
	}
	fmt.Printf("List of Events: %v\n", r)
}

func bookTickets(client ticket.TicketServiceClient, eventID string, numTickets int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.BookTickets(ctx, &ticket.BookRequest{
		EventId:    eventID,
		NumTickets: int32(numTickets),
	})
	if err != nil {
		log.Fatalf("Could not book tickets: %v", err)
	}
	fmt.Printf("Booked Tickets: %v\n", r)
}
