package handlers

import "net/http"

type TravelHandler interface {
	CreateTravel(w http.ResponseWriter, r *http.Request)
	SetTravelPreview(w http.ResponseWriter, r *http.Request)
	GetTravel(w http.ResponseWriter, r *http.Request)
	UpdateTravel(w http.ResponseWriter, r *http.Request)
	DeleteTravel(w http.ResponseWriter, r *http.Request)
}

type PlaceHandler interface {
	CreatePlace(w http.ResponseWriter, r *http.Request)
	SetPreview(w http.ResponseWriter, r *http.Request)
	SetImages(w http.ResponseWriter, r *http.Request)
	DeletePlace(w http.ResponseWriter, r *http.Request)
	UpdatePlace(w http.ResponseWriter, r *http.Request)
}

type ExpensesHandler interface {
	CreateExpense(w http.ResponseWriter, r *http.Request)
	GetExpense(w http.ResponseWriter, r *http.Request)
	UpdateExpense(w http.ResponseWriter, r *http.Request)
	DeleteExpense(w http.ResponseWriter, r *http.Request)
}
