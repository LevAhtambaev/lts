package handlers

import "net/http"

type TravelHandler interface {
	CreateTravel(w http.ResponseWriter, r *http.Request)
	SetTravelPreview(w http.ResponseWriter, r *http.Request)
}

type ExpenseHandler interface {
	CreateExpense(w http.ResponseWriter, r *http.Request)
}
