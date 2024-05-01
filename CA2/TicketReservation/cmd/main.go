package cmd

import (
	"TicketReservation/internal"
	Distrubted_Systems_git "TicketReservation/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	Distrubted_Systems_git.RegisterTicketServiceServer(s, &internal.TicketService{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
