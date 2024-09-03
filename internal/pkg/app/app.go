package app

import (
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	_ "lts/docs"
	"lts/internal/app/config"
	"lts/internal/app/handlers"
	"lts/internal/app/middleware"
	"lts/internal/app/repository"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	postgres = "postgres"
	sslMode  = "disable"
)

type App struct {
	ctx    context.Context
	cfg    config.Config
	logger *zap.SugaredLogger
}

func New(ctx context.Context, cfg config.Config, logger *zap.SugaredLogger) *App {
	return &App{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run() error {
	a.logger.Info("[app.Run]: the application is running")
	err := a.StartServer()
	if err != nil {
		return err
	}
	a.logger.Info("[app.Run]: the application is shut down")

	return nil
}

func (a *App) StartServer() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", a.cfg.PostgresConfig.Host, a.cfg.PostgresConfig.Port, a.cfg.PostgresConfig.User, a.cfg.PostgresConfig.Password, a.cfg.PostgresConfig.Name, sslMode)

	db, err := sqlx.Connect(postgres, dsn)
	if err != nil {
		return fmt.Errorf("[sqlx.Connect]: %w", err)
	}

	travelRepo := repository.NewTravelRepo(db)
	placeRepo := repository.NewPlaceRepositoryImpl(db)
	expenseRepo := repository.NewExpensesRepo(db)

	travelHandler := handlers.NewTravelHandlerImpl(travelRepo, placeRepo, expenseRepo, a.logger)
	th := handlers.TravelHandlerImplemented{TravelHandler: travelHandler}

	placesHandler := handlers.NewPlaceHandlerImpl(placeRepo, travelRepo, a.logger)
	ph := handlers.PlaceHandlerImplemented{PlaceHandler: placesHandler}

	expensesHandler := handlers.NewExpensesHandlerImpl(expenseRepo, placeRepo, a.logger)
	eh := handlers.ExpensesHandlerImplemented{ExpensesHandler: expensesHandler}

	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/travel", th.CreateTravel).Methods("POST", "OPTIONS")
	api.HandleFunc("/travel/{uuid}", th.GetTravel).Methods("GET")
	api.HandleFunc("/travel/preview/{uuid}", th.SetTravelPreview).Methods("PUT")
	api.HandleFunc("/travel/{uuid}", th.UpdateTravel).Methods("PUT")
	api.HandleFunc("/travel/{uuid}", th.DeleteTravel).Methods("DELETE")

	api.HandleFunc("/place/{travel_uuid}", ph.CreatePlace).Methods("POST")
	api.HandleFunc("/place/{travel_uuid}/{place_uuid}", ph.SetPreview).Methods("PUT")
	api.HandleFunc("/place/{travel_uuid}/{place_uuid}", ph.DeletePlace).Methods("DELETE")
	api.HandleFunc("/place/images/{travel_uuid}/{place_uuid}", ph.SetImages).Methods("PUT")
	api.HandleFunc("/place/{uuid}", ph.UpdatePlace).Methods("PUT")

	api.HandleFunc("/expenses/{place_uuid}", eh.CreateExpense).Methods("POST")
	api.HandleFunc("/expenses/{uuid}", eh.GetExpense).Methods("GET")
	api.HandleFunc("/expenses/{uuid}", eh.UpdateExpense).Methods("PUT")
	api.HandleFunc("/expenses/{uuid}", eh.DeleteExpense).Methods("DELETE")

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler).Methods("GET")

	router := middleware.LogMiddleware(a.logger, r)

	log.Println("server started")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		a.logger.Fatal()
	}

	return nil
}
