package main

import (
	"TicketReservation/internal"
	Distrubted_Systems_git "TicketReservation/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8998")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ticketService := internal.NewTicketService()
	Distrubted_Systems_git.RegisterTicketServiceServer(s, ticketService)
	log.Print("Starting gRPC Server ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
