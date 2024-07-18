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

type ExpensesHandlerImplemented struct {
	ExpensesHandler
}

type ExpensesHandlerImpl struct {
	ExpensesRepo repository.ExpensesRepository
	PlaceRepo    repository.PlaceRepository
	Logger       *zap.SugaredLogger
}

func NewExpensesHandlerImpl(expensesRepo repository.ExpensesRepository, placeRepo repository.PlaceRepository, logger *zap.SugaredLogger) *ExpensesHandlerImpl {
	return &ExpensesHandlerImpl{ExpensesRepo: expensesRepo, PlaceRepo: placeRepo, Logger: logger}
}

// CreateExpense godoc
// @Summary      Create a new expense
// @Description  Create a new expense entry for a specific place
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Param        place_uuid path string true "UUID of the place"
// @Param        expense body ds.Expense true "Expense details"
// @Success      201 {object} ds.Expense "Successfully created expense"
// @Failure      400 "Invalid place UUID or expense data"
// @Failure      500 "Internal server error"
// @Router       /expenses/{place_uuid} [post]
func (eh ExpensesHandlerImpl) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense ds.Expense

	vars := mux.Vars(r)
	uuidStr, ok := vars["place_uuid"]
	if !ok {
		eh.Logger.Info("uuid is missing in parameters")
	}

	uuidParsed, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	expense, err = eh.ExpensesRepo.CreateExpense(r.Context(), expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = eh.PlaceRepo.SetExpenses(r.Context(), expense.ID, uuidParsed)
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

// GetExpense godoc
// @Summary      Get expense details
// @Description  Retrieve details of a specific expense by its UUID
// @Tags         Expenses
// @Produce      json
// @Param        uuid path string true "UUID of the expense"
// @Success      200 {object} ds.Expense "Successfully retrieved expense details"
// @Failure      400 "Invalid UUID format"
// @Failure      500 "Internal server error"
// @Router       /expenses/{uuid} [get]
func (eh ExpensesHandlerImpl) GetExpense(w http.ResponseWriter, r *http.Request) {
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

	expense, err := eh.ExpensesRepo.GetExpense(r.Context(), uuidParsed)
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

// UpdateExpense godoc
// @Summary      Update expense details
// @Description  Update the details of a specific expense by its UUID
// @Tags         Expenses
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID of the expense"
// @Param        expense body ds.Expense true "Expense details"
// @Success      200 {object} ds.Expense "Successfully updated expense details"
// @Failure      400 "Invalid UUID format or invalid expense data"
// @Failure      500 "Internal server error"
// @Router       /expenses/{uuid} [put]
func (eh ExpensesHandlerImpl) UpdateExpense(w http.ResponseWriter, r *http.Request) {
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

	err = json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	expense, err = eh.ExpensesRepo.UpdateExpense(r.Context(), expense, uuidParsed)
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

// DeleteExpense godoc
// @Summary      Delete an expense
// @Description  Delete a specific expense by its UUID
// @Tags         Expenses
// @Produce      json
// @Param        uuid path string true "UUID of the expense"
// @Success      200 "Successfully deleted expense"
// @Failure      400 "Invalid UUID format"
// @Failure      500 "Internal server error"
// @Router       /expenses/{uuid} [delete]
func (eh ExpensesHandlerImpl) DeleteExpense(w http.ResponseWriter, r *http.Request) {
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

	err = eh.ExpensesRepo.DeleteExpense(r.Context(), uuidParsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
