package usecases

import "github.com/gabrielmq/events-api/internal/events/domain"

type ListSpotsUseCase struct {
	repository domain.EventRepository
}

func NewListSpotsUseCase(repository domain.EventRepository) *ListSpotsUseCase {
	return &ListSpotsUseCase{repository: repository}
}

func (u ListSpotsUseCase) Execute(eventId string) (*ListSpotsOutput, error) {
	event, err := u.repository.FindEventByID(eventId)
	if err != nil {
		return nil, err
	}

	spots, err := u.repository.FindSpotsByEventID(eventId)
	if err != nil {
		return nil, err
	}

	spotsOutput := make([]SpotOutput, len(spots))
	for i, spot := range spots {
		spotsOutput[i] = SpotOutput{
			ID:       spot.ID,
			Name:     spot.Name,
			Status:   string(spot.Status),
			TicketID: spot.TicketID,
		}
	}

	eventOutput := EventOutput{
		ID:           event.ID,
		Name:         event.Name,
		Location:     event.Location,
		Organization: event.Organization,
		Rating:       string(event.Rating),
		Date:         event.Date.Format("2006-01-2 15:04:05"),
		ImageURL:     event.ImageURL,
		Capacity:     event.Capacity,
		Price:        event.Price,
		PartnerID:    event.PartnerID,
	}

	return &ListSpotsOutput{
		Event: eventOutput,
		Spots: spotsOutput,
	}, nil
}

type ListSpotsOutput struct {
	Event EventOutput  `json:"event"`
	Spots []SpotOutput `json:"spots"`
}
