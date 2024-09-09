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
	TravelRepo   repository.TravelRepository
	PlaceRepo    repository.PlaceRepository
	ExpensesRepo repository.ExpensesRepository
	Logger       *zap.SugaredLogger
}

func NewTravelHandlerImpl(travelRepo repository.TravelRepository, placeRepo repository.PlaceRepository, expensesRepo repository.ExpensesRepository, logger *zap.SugaredLogger) *TravelHandlerImpl {
	return &TravelHandlerImpl{
		TravelRepo:   travelRepo,
		PlaceRepo:    placeRepo,
		ExpensesRepo: expensesRepo,
		Logger:       logger,
	}
}

// CreateTravel godoc
// @Summary      Create a new travel
// @Description  Create a new travel entry with provided details
// @Tags         Travel
// @Accept       json
// @Produce      json
// @Param        travel body ds.Travel true "Travel details"
// @Success      201 {object} ds.Travel "Successfully created travel"
// @Failure      400 "Invalid travel data"
// @Failure      500 "Internal server error"
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

// SetTravelPreview godoc
// @Summary      Set a preview for travel
// @Description  Set a preview picture for travel
// @Tags         Travel
// @Accept       multipart/form-data
// @Produce      json
// @Param        uuid path string true "UUID of the travel"
// @Param        file formData file true "Preview picture"
// @Success      200 "Successfully set preview"
// @Failure      400 "Invalid UUID format"
// @Failure      500 "Internal server error"
// @Router       /travel/preview/{uuid} [put]
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

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(file)

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

// GetTravel godoc
// @Summary      Get travel details
// @Description  Retrieve detailed information about specific travel including places and images
// @Tags         Travel
// @Produce      json
// @Param        uuid path string true "UUID of the travel"
// @Success      200 {object} ds.FullTravel "Successfully retrieved travel details"
// @Failure      400 "Invalid UUID format"
// @Failure      500 "Internal server error"
// @Router       /travel/{uuid} [get]
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

	if travel.Preview != "" {
		travel.Preview, err = helpers.LoadImage(travel.Preview)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var places []ds.FullPlace

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
			if imagePath != "" {
				imageData, err := helpers.LoadImage(imagePath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				place.Images[i] = imageData
			}

		}
		th.Logger.Info(place.Preview)
		if place.Preview != "" {
			place.Preview, err = helpers.LoadImage(place.Preview)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var expense ds.Expense

		if place.Expenses != uuid.Nil {
			expense, err = th.ExpensesRepo.GetExpense(r.Context(), place.Expenses)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		fullPlace := ds.FullPlace{
			ID:      place.ID,
			Name:    place.Name,
			Story:   place.Story,
			Date:    place.Date,
			Images:  place.Images,
			Preview: place.Preview,
		}

		// Записываем expense в Expenses только если place.Expenses не nil
		if place.Expenses != uuid.Nil {
			fullPlace.Expenses = expense
		}

		places = append(places, fullPlace)
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

// UpdateTravel godoc
// @Summary      Update travel details
// @Description  Update the details of a specific travel
// @Tags         Travel
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID of the travel"
// @Param        travel body ds.Travel true "Travel details to update"
// @Success      200 "Successfully updated travel details"
// @Failure      400 "Invalid UUID format or invalid travel data"
// @Failure      500 "Internal server error"
// @Router       /travel/{uuid} [put]
func (th *TravelHandlerImpl) UpdateTravel(w http.ResponseWriter, r *http.Request) {
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

	err = json.NewDecoder(r.Body).Decode(&travel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = th.TravelRepo.UpdateTravel(r.Context(), UUID, travel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteTravel godoc
// @Summary      Delete travel
// @Description  Delete a specific travel and all associated places and expenses
// @Tags         Travel
// @Produce      json
// @Param        uuid path string true "UUID of the travel"
// @Success      200 "Successfully deleted travel"
// @Failure      400 "Invalid UUID format"
// @Failure      500 "Internal server error"
// @Router       /travel/{uuid} [delete]
func (th *TravelHandlerImpl) DeleteTravel(w http.ResponseWriter, r *http.Request) {
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

	err = th.TravelRepo.DeleteTravel(r.Context(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, placeUUID := range travel.Places {
		place, err := th.PlaceRepo.GetPlace(r.Context(), placeUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = th.ExpensesRepo.DeleteExpense(r.Context(), place.Expenses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = th.PlaceRepo.DeletePlace(r.Context(), placeUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	path := fmt.Sprintf("./images/travel/%s", UUID)

	err = os.RemoveAll(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
