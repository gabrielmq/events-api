package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/gabrielmq/events-api/docs"
	"github.com/gabrielmq/events-api/internal/events/application/events/usecases"
	httpHandlres "github.com/gabrielmq/events-api/internal/events/infraestructure/http"
	"github.com/gabrielmq/events-api/internal/events/infraestructure/repository"
	"github.com/gabrielmq/events-api/internal/events/infraestructure/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Events API
// @version 1.0
// @description This is a server for managing events. Imersão Full Cycle
// @host localhost:8080
// @BasePath /
func main() {
	db, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/events")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventRepepository, err := repository.NewMysqlEventRepository(db)
	if err != nil {
		panic(err)
	}

	partinerBaseURLs := map[int]string{
		1: "http://localhost:9080/partner1",
		2: "http://localhost:9080/partner2",
	}

	partnerFactory := service.NewPartnerFactory(partinerBaseURLs)

	listEventsUseCase := usecases.NewListEventsUseCase(eventRepepository)
	getEventUseCase := usecases.NewGetEventUseCase(eventRepepository)
	createEventUseCase := usecases.NewCreateEventUseCase(eventRepepository)
	buyTicketsUseCase := usecases.NewBuyTicketsUseCase(eventRepepository, partnerFactory)
	createSpotsUseCase := usecases.NewCreateSpotsUseCase(eventRepepository)
	listSpotsUseCase := usecases.NewListSpotsUseCase(eventRepepository)

	eventsHandler := httpHandlres.NewEventsHandler(
		listEventsUseCase,
		getEventUseCase,
		buyTicketsUseCase,
		listSpotsUseCase,
		createSpotsUseCase,
		createEventUseCase,
	)

	r := http.NewServeMux()
	r.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	r.HandleFunc("/events", eventsHandler.ListEvents)
	r.HandleFunc("/events/{eventID}", eventsHandler.GetEvent)
	r.HandleFunc("/events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /events", eventsHandler.CreateEvent)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)
	r.HandleFunc("POST /events/{eventID}/spots", eventsHandler.CreateSpots)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Println("Iniciando graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Erro durante o graceful shutdown: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("Iniciando servidor HTTP na porta 8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor: %v\n", err)
	}

	<-idleConnsClosed
	log.Println("Servidor HTTP finalizado")
}
