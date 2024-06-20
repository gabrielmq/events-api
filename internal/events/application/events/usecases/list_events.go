package usecases

import "github.com/gabrielmq/events-api/internal/events/domain"

type ListEventsUseCase struct {
	repository domain.EventRepository
}

func NewListEventsUseCase(repository domain.EventRepository) *ListEventsUseCase {
	return &ListEventsUseCase{repository: repository}
}

func (u ListEventsUseCase) Execute() (*ListEventsOutput, error) {
	events, err := u.repository.ListEvents()
	if err != nil {
		return nil, err
	}

	eventsOutput := make([]EventOutput, len(events))
	for i, event := range events {
		eventsOutput[i] = EventOutput{
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
	}

	return &ListEventsOutput{Events: eventsOutput}, nil
}

type ListEventsOutput struct {
	Events []EventOutput `json:"events"`
}
