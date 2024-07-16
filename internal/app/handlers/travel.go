package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"lts/internal/app/ds"
	"lts/internal/app/repository"
	"net/http"
	"os"
)

type TravelHandlerImpl struct {
	TravelRepo repository.TravelRepository
	Logger     *zap.SugaredLogger
}

func NewTravelHandlerImpl(travelRepo repository.TravelRepository, logger *zap.SugaredLogger) *TravelHandlerImpl {
	return &TravelHandlerImpl{
		TravelRepo: travelRepo,
		Logger:     logger,
	}
}

func (th *TravelHandlerImpl) CreateTravel(w http.ResponseWriter, r *http.Request) {
	var travel ds.Travel

	err := json.NewDecoder(r.Body).Decode(&travel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	travel, err = th.TravelRepo.CreateTravel(r.Context(), travel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(travel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (th *TravelHandlerImpl) SetTravelPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuidStr, ok := vars["uuid"]
	if !ok {
		th.Logger.Info("uuid is missing in parameters")
	}

	uuidParsed, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	th.Logger.Info(uuidStr)

	err = os.MkdirAll(fmt.Sprintf("./images/travel/%s", uuidStr), os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("./images/travel/%s/preview.jpg", uuidStr)

	file, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = th.TravelRepo.SetTravelPreview(r.Context(), path, uuidParsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
