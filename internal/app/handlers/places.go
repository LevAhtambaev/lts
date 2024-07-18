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

type PlaceHandlerImplemented struct {
	PlaceHandler
}

type PlaceHandlerImpl struct {
	PlaceRepo  repository.PlaceRepository
	TravelRepo repository.TravelRepository
	Logger     *zap.SugaredLogger
}

func NewPlaceHandlerImpl(placeRepo repository.PlaceRepository, travelRepo repository.TravelRepository, logger *zap.SugaredLogger) *PlaceHandlerImpl {
	return &PlaceHandlerImpl{PlaceRepo: placeRepo, TravelRepo: travelRepo, Logger: logger}
}

// CreatePlace godoc
// @Summary      Create a new place
// @Description  Create a new place and associate it with a specific travel
// @Tags         Places
// @Accept       json
// @Produce      json
// @Param        travel_uuid path string true "UUID of the travel"
// @Param        place body ds.Place true "Place details"
// @Success      201 {object} ds.Place "Successfully created place"
// @Failure      400 "Invalid travel UUID or place data"
// @Failure      500 "Internal server error"
// @Router       /place/{travel_uuid} [post]
func (ph PlaceHandlerImpl) CreatePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	travelStr, ok := vars["travel_uuid"]
	if !ok {
		ph.Logger.Info("travel uuid is missing in parameters")
	}

	travelUUID, err := uuid.Parse(travelStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var place ds.Place

	err = json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	place, err = ph.PlaceRepo.CreatePlace(r.Context(), place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ph.TravelRepo.AddPlace(r.Context(), travelUUID, place.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SetPreview godoc
// @Summary      Set a preview for a place
// @Description  Set a preview picture for a place
// @Tags         Places
// @Accept       multipart/form-data
// @Produce      json
// @Param        travel_uuid path string true "UUID of the travel"
// @Param        place_uuid path string true "UUID of the place"
// @Param        file formData file true "Preview picture"
// @Success      200 "Successfully set preview"
// @Failure      400 "Invalid travel UUID or place UUID"
// @Failure      500 "Internal server error"
// @Router       /place/{travel_uuid}/{place_uuid} [put]
func (ph PlaceHandlerImpl) SetPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	travelStr, ok := vars["travel_uuid"]
	if !ok {
		ph.Logger.Info("travel uuid is missing in parameters")
	}

	travelUUID, err := uuid.Parse(travelStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	placeStr, ok := vars["place_uuid"]
	if !ok {
		ph.Logger.Info("place uuid is missing in parameters")
	}

	placeUUID, err := uuid.Parse(placeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = os.MkdirAll(fmt.Sprintf("./images/travel/%s/places/%s", travelUUID, placeUUID), os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("./images/travel/%s/places/%s/preview.jpg", travelUUID, placeUUID)

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

	err = ph.PlaceRepo.SetPreview(r.Context(), path, placeUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// SetImages godoc
// @Summary      Set images for a place
// @Description  Upload images for a specific place associated with a travel
// @Tags         Places
// @Accept       multipart/form-data
// @Produce      json
// @Param        travel_uuid path string true "UUID of the travel"
// @Param        place_uuid path string true "UUID of the place"
// @Param        image formData file true "Image file"
// @Success      200 "Successfully set images"
// @Failure      400 "Invalid travel UUID or place UUID"
// @Failure      500 "Internal server error"
// @Router       /place/images/{travel_uuid}/{place_uuid} [put]
func (ph PlaceHandlerImpl) SetImages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	travelStr, ok := vars["travel_uuid"]
	if !ok {
		ph.Logger.Info("travel uuid is missing in parameters")
	}

	travelUUID, err := uuid.Parse(travelStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	placeStr, ok := vars["place_uuid"]
	if !ok {
		ph.Logger.Info("place uuid is missing in parameters")
	}

	placeUUID, err := uuid.Parse(placeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var paths []string

	ph.Logger.Info(r.MultipartForm)

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка при парсинге формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil {
		http.Error(w, "Ошибка: форма не содержит данных", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["image"]
	if files == nil {
		http.Error(w, "Ошибка: не найдены файлы с ключом 'image'", http.StatusBadRequest)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = os.MkdirAll(fmt.Sprintf("./images/travel/%s/places/%s/images", travelUUID, placeUUID), os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		path := fmt.Sprintf("./images/travel/%s/places/%s/images/%s", travelUUID, placeUUID, fileHeader.Filename)

		paths = append(paths, path)

		out, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = ph.PlaceRepo.SetImages(r.Context(), paths, placeUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeletePlace godoc
// @Summary      Delete a place
// @Description  Delete a specific place associated with a travel, including all associated data and images
// @Tags         Places
// @Produce      json
// @Param        travel_uuid path string true "UUID of the travel"
// @Param        place_uuid path string true "UUID of the place"
// @Success      200 "Successfully deleted place"
// @Failure      400 "Invalid travel UUID or place UUID"
// @Failure      500 "Internal server error"
// @Router       /place/{travel_uuid}/{place_uuid} [delete]
func (ph PlaceHandlerImpl) DeletePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	travelStr, ok := vars["travel_uuid"]
	if !ok {
		ph.Logger.Info("travel uuid is missing in parameters")
	}

	travelUUID, err := uuid.Parse(travelStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	placeStr, ok := vars["place_uuid"]
	if !ok {
		ph.Logger.Info("place uuid is missing in parameters")
	}

	placeUUID, err := uuid.Parse(placeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("./images/travel/%s/places/%s", travelUUID, placeUUID)

	err = os.RemoveAll(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ph.PlaceRepo.DeletePlace(r.Context(), placeUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdatePlace godoc
// @Summary      Update place details
// @Description  Update the details of a specific place associated with a travel
// @Tags         Places
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID of the place"
// @Param        place body ds.Place true "Place details"
// @Success      200 "Successfully updated place details"
// @Failure      400 "Invalid UUID format or invalid place data"
// @Failure      500 "Internal server error"
// @Router       /place/{uuid} [put]
func (ph PlaceHandlerImpl) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuidStr, ok := vars["uuid"]
	if !ok {
		ph.Logger.Info("uuid is missing in parameters")
	}

	UUID, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var place ds.Place

	err = json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ph.PlaceRepo.UpdatePlace(r.Context(), UUID, place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
