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
	travelHandler := handlers.NewTravelHandlerImpl(travelRepo, a.logger)
	th := handlers.TravelHandlerImplemented{TravelHandler: travelHandler}

	expenseRepo := repository.NewExpensesRepo(db)
	expensesHandler := handlers.NewExpensesHandlerImpl(expenseRepo, a.logger)
	eh := handlers.ExpensesHandlerImplemented{ExpensesHandler: expensesHandler}

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/travel", th.CreateTravel).Methods("POST")
	api.HandleFunc("/travel/{uuid}", th.SetTravelPreview).Methods("PUT")

	api.HandleFunc("/expenses", eh.CreateExpense).Methods("POST")
	api.HandleFunc("/expenses/{uuid}", eh.GetExpense).Methods("GET")
	api.HandleFunc("/expenses/{uuid}", eh.UpdateExpense).Methods("PUT")
	api.HandleFunc("/expenses/{uuid}", eh.DeleteExpense).Methods("DELETE")

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler).Methods("GET")

	router := middleware.LogMiddleware(a.logger, r)

	log.Println("server started")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		a.logger.Fatal()
	}

	return nil
}
