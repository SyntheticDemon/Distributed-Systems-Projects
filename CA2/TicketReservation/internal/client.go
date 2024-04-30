package internal

import "TicketReservation/pkg/models"

type TicketClient interface {
	ListAvailableEvents() ([]*models.Event, error)
	BookTicket(eventID string, numTickets int) ([]string, error)
}

type TicketReservationClient struct {
	service *TicketService
}

func (trc *TicketReservationClient) BookTickets(eventID string, numTickets int) ([]string, error) {
	return trc.service.BookTickets(eventID, numTickets)
}

func (trc *TicketReservationClient) ListAvailableEvents() []*models.Event {
	return trc.service.ListEvents()
}
