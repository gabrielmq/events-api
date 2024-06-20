package usecases

import (
	"fmt"

	"github.com/gabrielmq/events-api/internal/events/domain"
)

type CreateSpotsUseCase struct {
	repo domain.EventRepository
}

func NewCreateSpotsUseCase(repo domain.EventRepository) *CreateSpotsUseCase {
	return &CreateSpotsUseCase{repo: repo}
}

func (uc *CreateSpotsUseCase) Execute(input CreateSpotsInput) (*CreateSpotsOutput, error) {
	event, err := uc.repo.FindEventByID(input.EventID)
	if err != nil {
		return nil, err
	}

	spots := make([]domain.Spot, input.NumberOfSpots)
	for i := 0; i < input.NumberOfSpots; i++ {
		spotName := generateSpotName(i)
		spot, err := domain.NewSpot(event, spotName)
		if err != nil {
			return nil, err
		}
		if err := uc.repo.CreateSpot(spot); err != nil {
			return nil, err
		}
		spots[i] = *spot
	}

	spotDTOs := make([]SpotOutput, len(spots))
	for i, spot := range spots {
		spotDTOs[i] = SpotOutput{
			ID:       spot.ID,
			Name:     spot.Name,
			Status:   string(spot.Status),
			TicketID: spot.TicketID,
		}
	}

	return &CreateSpotsOutput{Spots: spotDTOs}, nil
}

func generateSpotName(index int) string {
	// Gera um nome de spot baseado no índice. Ex: A1, A2, ..., B1, B2, etc.
	letter := 'A' + rune(index/10)
	number := index%10 + 1
	return fmt.Sprintf("%c%d", letter, number)
}

type CreateSpotsInput struct {
	EventID       string `json:"event_id"`
	NumberOfSpots int    `json:"number_of_spots"`
}

type CreateSpotsOutput struct {
	Spots []SpotOutput `json:"spots"`
}
