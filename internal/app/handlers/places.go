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
