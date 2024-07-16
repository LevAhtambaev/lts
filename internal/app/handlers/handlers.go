package handlers

import "net/http"

type TravelHandler interface {
	CreateTravel(w http.ResponseWriter, r *http.Request)
	SetTravelPreview(w http.ResponseWriter, r *http.Request)
}

type ExpensesHandler interface {
	CreateExpense(w http.ResponseWriter, r *http.Request)
	GetExpense(w http.ResponseWriter, r *http.Request)
	UpdateExpense(w http.ResponseWriter, r *http.Request)
	DeleteExpense(w http.ResponseWriter, r *http.Request)
}
