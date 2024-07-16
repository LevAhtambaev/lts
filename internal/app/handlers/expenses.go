package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"lts/internal/app/ds"
	"lts/internal/app/repository"
	"net/http"
)

type ExpensesHandlerImpl struct {
	ExpensesRepo repository.ExpensesRepository
	Logger       *zap.SugaredLogger
}

func NewExpensesHandlerImpl(expensesRepo repository.ExpensesRepository, logger *zap.SugaredLogger) *ExpensesHandlerImpl {
	return &ExpensesHandlerImpl{ExpensesRepo: expensesRepo, Logger: logger}
}

func (eh ExpensesHandlerImpl) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense ds.Expense

	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	expense, err = eh.ExpensesRepo.CreateExpense(r.Context(), expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (eh ExpensesHandlerImpl) GetExpense(w http.ResponseWriter, r *http.Request) {
	var expense ds.Expense

	vars := mux.Vars(r)
	uuidStr, ok := vars["uuid"]
	if !ok {
		eh.Logger.Info("uuid is missing in parameters")
	}

	uuidParsed, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expense, err = eh.ExpensesRepo.GetExpense(r.Context(), uuidParsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
