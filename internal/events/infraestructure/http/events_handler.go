package http

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielmq/events-api/internal/events/application/events/usecases"
)

type EventsHandler struct {
	listEventsUseCase  *usecases.ListEventsUseCase
	getEventUseCase    *usecases.GetEventUseCase
	buyTicketsUseCase  *usecases.BuyTicketsUseCase
	listSpotsUseCase   *usecases.ListSpotsUseCase
	createSpotsUseCase *usecases.CreateSpotsUseCase
	createEventUseCase *usecases.CreateEventUseCase
}

func NewEventsHandler(
	listEventsUseCase *usecases.ListEventsUseCase,
	getEventUseCase *usecases.GetEventUseCase,
	buyTicketsUseCase *usecases.BuyTicketsUseCase,
	listSpotsUseCase *usecases.ListSpotsUseCase,
	createSpotsUseCase *usecases.CreateSpotsUseCase,
	createEventUseCase *usecases.CreateEventUseCase,
) *EventsHandler {
	return &EventsHandler{
		listEventsUseCase:  listEventsUseCase,
		getEventUseCase:    getEventUseCase,
		buyTicketsUseCase:  buyTicketsUseCase,
		listSpotsUseCase:   listSpotsUseCase,
		createSpotsUseCase: createSpotsUseCase,
		createEventUseCase: createEventUseCase,
	}
}

// ListEvents handles the request to list all events.
// @Summary List all events
// @Description Get all events with their details
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} usecases.ListEventsOutput
// @Failure 500 {object} string
// @Router /events [get]
func (h *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	output, err := h.listEventsUseCase.Execute()
	if err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// GetEvent handles the request to get details of a specific event.
// @Summary Get event details
// @Description Get details of an event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Param eventID path string true "Event ID"
// @Success 200 {object} usecases.GetEventOutput
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /events/{eventID} [get]
func (h *EventsHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	output, err := h.getEventUseCase.Execute(eventID)
	if err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// BuyTickets handles the request to buy tickets for an event.
// @Summary Buy tickets for an event
// @Description Buy tickets for a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Param input body usecases.BuyTicketsInput true "Input data"
// @Success 200 {object} usecases.BuyTicketsOutput
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /checkout [post]
func (h *EventsHandler) BuyTickets(w http.ResponseWriter, r *http.Request) {
	var input usecases.BuyTicketsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.buyTicketsUseCase.Execute(input)
	if err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// ListSpots lists spots for an event.
// @Summary List spots for an event
// @Description List all spots for a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Param eventID path string true "Event ID"
// @Success 200 {object} usecases.ListSpotsOutput
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /events/{eventID}/spots [get]
func (h *EventsHandler) ListSpots(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	output, err := h.listSpotsUseCase.Execute(eventID)
	if err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// CreateEvent handles the request to create a new event.
// @Summary Create a new event
// @Description Create a new event with the given details
// @Tags Events
// @Accept json
// @Produce json
// @Param input body usecases.CreateEventInput true "Input data"
// @Success 201 {object} usecases.CreateEventOutput
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /events [post]
func (h *EventsHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input usecases.CreateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.createEventUseCase.Execute(input)
	if err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// CreateSpots handles the creation of spots.
// @Summary Create spots for an event
// @Description Create a specified number of spots for an event
// @Tags Events
// @Accept json
// @Produce json
// @Param eventID path string true "Event ID"
// @Param input body CreateSpotsRequest true "Input data"
// @Success 201 {object} usecases.CreateSpotsOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events/{eventID}/spots [post]
func (h *EventsHandler) CreateSpots(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	var input usecases.CreateSpotsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.EventID = eventID

	output, err := h.createSpotsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EventsHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateSpotsRequest struct {
	NumberOfSpots int `json:"number_of_spots"`
}
