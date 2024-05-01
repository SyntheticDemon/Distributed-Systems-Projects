package main

import (
	proto "TicketReservation/proto"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: <command> [arguments]")
		return
	}
	command := os.Args[1]

	conn, err := grpc.Dial("localhost:8998", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	client := proto.NewTicketServiceClient(conn)
	ctx := context.Background()

	switch command {
	case "create":
		createEvent(ctx, client, os.Args[2:])
	case "book":
		bookTickets(ctx, client, os.Args[2:])
	case "list":
		listEvents(ctx, client)
	case "stress":
		stressTest(ctx, client, os.Args[2:])
	default:
		fmt.Println("Unknown command")
	}
}

func createEvent(ctx context.Context, client proto.TicketServiceClient, args []string) {
	if len(args) != 3 {
		fmt.Println("usage: create <name> <date> <totalTickets>")
		return
	}
	name := args[0]
	date, _ := strconv.ParseInt(args[1], 10, 64)
	totalTickets, _ := strconv.Atoi(args[2])
	req := &proto.CreateEventRequest{
		Name:         name,
		Date:         date,
		TotalTickets: int32(totalTickets),
	}
	resp, err := client.CreateEvent(ctx, req)
	if err != nil {
		fmt.Printf("Error creating event: %v\n", err)
		return
	}
	fmt.Printf("Event created: ID %s, Name %s\n", resp.Id, resp.Name)
}

func bookTickets(ctx context.Context, client proto.TicketServiceClient, args []string) {
	if len(args) != 2 {
		fmt.Println("usage: book <eventId> <numTickets>")
		return
	}
	eventId := args[0]
	numTickets, _ := strconv.Atoi(args[1])
	req := &proto.BookRequest{
		EventId:    eventId,
		NumTickets: int32(numTickets),
	}
	resp, err := client.BookTickets(ctx, req)
	if err != nil {
		fmt.Printf("Error booking tickets: %v\n", err)
		return
	}
	fmt.Printf("Tickets booked: %v\n", resp.TicketIds)
}

func listEvents(ctx context.Context, client proto.TicketServiceClient) {
	resp, err := client.ListEvents(ctx, &proto.ListEventsRequest{})
	if err != nil {
		fmt.Printf("Error listing events: %v\n", err)
		return
	}
	for _, event := range resp.EventDetails {
		fmt.Printf("Event ID: %s, Name: %s, Date: %d, TotalTickets: %d, AvailableTickets: %d\n",
			event.Id, event.Name, event.Date, event.TotalTickets, event.AvailableTickets)
	}
}

func stressTest(ctx context.Context, client proto.TicketServiceClient, args []string) {
	if len(args) != 4 {
		fmt.Println("usage: stress <eventId> <numTickets> <count> <interval>")
		return
	}
	eventId := args[0]
	numTickets, _ := strconv.Atoi(args[1])
	count, _ := strconv.Atoi(args[2])
	interval, _ := strconv.ParseInt(args[3], 10, 64)

	req := &proto.BookRequest{
		EventId:    eventId,
		NumTickets: int32(numTickets),
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()
	for i := 0; i < count; i++ {
		<-ticker.C
		_, err := client.BookTickets(ctx, req)
		if err != nil {
			fmt.Printf("Error at request %d: %v\n", i+1, err)
		} else {
			fmt.Printf("Request %d successful\n", i+1)
		}
	}
}
