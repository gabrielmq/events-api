package usecases

import (
	"github.com/gabrielmq/events-api/internal/events/domain"
	"github.com/gabrielmq/events-api/internal/events/infraestructure/service"
)

type BuyTicketsUseCase struct {
	repository     domain.EventRepository
	partnerFactory service.PartnerFactory
}

func NewBuyTicketsUseCase(
	repository domain.EventRepository,
	partnerFactory service.PartnerFactory,
) *BuyTicketsUseCase {
	return &BuyTicketsUseCase{repository: repository, partnerFactory: partnerFactory}
}

func (u BuyTicketsUseCase) Execute(input BuyTicketsInput) (*BuyTicketsOutput, error) {
	event, err := u.repository.FindEventByID(input.EventID)
	if err != nil {
		return nil, err
	}

	req := &service.ReservationRequest{
		EventID:    event.ID,
		Spots:      input.Spots,
		TicketType: input.TicketType,
		CardHash:   input.CardHash,
		Email:      input.Email,
	}

	partner, err := u.partnerFactory.CreatePartner(event.PartnerID)
	if err != nil {
		return nil, err
	}

	reservertions, err := partner.MakeReservation(req)
	if err != nil {
		return nil, err
	}

	tickets := make([]domain.Ticket, len(reservertions))
	for i, reservertion := range reservertions {
		spot, err := u.repository.FindSpotByName(event.ID, reservertion.Spot)
		if err != nil {
			return nil, err
		}

		ticket, err := domain.NewTicket(event, spot, domain.TicketType(input.TicketType))
		if err != nil {
			return nil, err
		}

		if err := u.repository.CreateTicket(ticket); err != nil {
			return nil, err
		}

		spot.Reserve(ticket.ID)
		if err := u.repository.ReserveSpot(spot.ID, ticket.ID); err != nil {
			return nil, err
		}

		tickets[i] = *ticket
	}

	ticketsOutput := make([]TicketOutput, len(tickets))
	for i, ticket := range tickets {
		ticketsOutput[i] = TicketOutput{
			ID:         ticket.ID,
			SpotID:     ticket.Spot.ID,
			TicketType: string(ticket.TicketType),
			Price:      ticket.Price,
		}
	}

	return &BuyTicketsOutput{Tickets: ticketsOutput}, nil
}

type BuyTicketsInput struct {
	EventID    string   `json:"event_id"`
	Spots      []string `json:"spots"`
	TicketType string   `json:"ticket_type"`
	CardHash   string   `json:"card_hash"`
	Email      string   `json:"email"`
}

type BuyTicketsOutput struct {
	Tickets []TicketOutput `json:"tickets"`
}
