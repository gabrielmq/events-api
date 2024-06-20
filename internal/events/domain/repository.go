package domain

type EventRepository interface {
	CreateEvent(event *Event) error
	CreateSpot(spot *Spot) error
	CreateTicket(ticket *Ticket) error
	ReserveSpot(spotId, ticketId string) error
	ListEvents() ([]Event, error)
	FindEventByID(eventId string) (*Event, error)
	FindSpotsByEventID(eventId string) ([]*Spot, error)
	FindSpotByName(eventId, spotName string) (*Spot, error)
}
