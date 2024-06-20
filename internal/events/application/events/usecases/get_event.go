package usecases

import "github.com/gabrielmq/events-api/internal/events/domain"

type GetEventUseCase struct {
	repository domain.EventRepository
}

func NewGetEventUseCase(repository domain.EventRepository) *GetEventUseCase {
	return &GetEventUseCase{repository: repository}
}

func (u GetEventUseCase) Execute(eventId string) (*GetEventOutput, error) {
	event, err := u.repository.FindEventByID(eventId)
	if err != nil {
		return nil, err
	}

	return &GetEventOutput{
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
	}, nil
}

type GetEventOutput struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Organization string  `json:"organization"`
	Rating       string  `json:"rating"`
	Date         string  `json:"date"`
	ImageURL     string  `json:"image_url"`
	Capacity     int     `json:"capacity"`
	Price        float64 `json:"price"`
	PartnerID    int     `json:"partner_id"`
}
