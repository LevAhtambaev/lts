package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"lts/internal/app/ds"
	"lts/internal/app/helpers"
	"lts/internal/app/repository"
	"net/http"
	"os"
)

type TravelHandlerImplemented struct {
	TravelHandler
}

type TravelHandlerImpl struct {
	TravelRepo repository.TravelRepository
	PlaceRepo  repository.PlaceRepository
	Logger     *zap.SugaredLogger
}

func NewTravelHandlerImpl(travelRepo repository.TravelRepository, placeRepo repository.PlaceRepository, logger *zap.SugaredLogger) *TravelHandlerImpl {
	return &TravelHandlerImpl{
		TravelRepo: travelRepo,
		PlaceRepo:  placeRepo,
		Logger:     logger,
	}
}

// CreateTravel godoc
// @Summary      Create travel
// @Description  Create new travel with description
// @Tags         Travel
// @Produce      json
// @Param travel body ds.Travel true "Create travel"
// @Success      201  {object}  ds.Travel
// @Failure 	 500
// @Router       /travel [post]
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

func (th *TravelHandlerImpl) GetTravel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuidStr, ok := vars["uuid"]
	if !ok {
		th.Logger.Info("uuid is missing in parameters")
	}

	UUID, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var travel ds.Travel

	travel, err = th.TravelRepo.GetTravel(r.Context(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	travel.Preview, err = helpers.LoadImage(travel.Preview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var places []ds.Place

	for _, placeUUID := range travel.Places {
		var place ds.Place
		place, err = th.PlaceRepo.GetPlace(r.Context(), placeUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images := place.Images
		for i, imagePath := range images {
			//if imagePath == "" {
			//	continue
			//}
			th.Logger.Info(imagePath)
			imageData, err := helpers.LoadImage(imagePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			place.Images[i] = imageData
		}
		th.Logger.Info(place.Preview)
		if place.Preview != "" {
			place.Preview, err = helpers.LoadImage(place.Preview)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		places = append(places, place)
	}

	fullTravel := ds.FullTravel{
		ID:          travel.ID,
		Name:        travel.Name,
		Description: travel.Description,
		DateStart:   travel.DateStart,
		DateEnd:     travel.DateEnd,
		Places:      places,
		Preview:     travel.Preview,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fullTravel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
