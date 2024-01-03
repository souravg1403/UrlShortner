package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"urlshortener/utils"
)

type URL struct {
	WebAddress string `json:"url"`
}

var PORT string = "5000"

func (h *Handler) SetURL(w http.ResponseWriter, r *http.Request) {
	var data URL
	err := json.NewDecoder(r.Body).Decode(&data)
	err = errors.New("Failed to decode body")
	if err != nil {
		// w = utils.JSONError(w, err, http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(data)

	ctx := context.Background()
	//sending the url to urlShortening for checking exsitence in db and set the value to db
	shortenedUrlHash, err := h.ShortenURL(ctx, data)
	if err != nil {
		w = utils.JSONError(w, err, http.StatusBadRequest)
		return
	}

	//after checking existence, setting into db
	_, err = h.KeyValueStore.SetKey(ctx, *shortenedUrlHash, data.WebAddress)
	if err != nil {
		w = utils.JSONError(w, err, http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"shortenedUrl": fmt.Sprintf("http://localhost:%s/%s", PORT, *shortenedUrlHash),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w = utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBytes)
}

func (h *Handler) FetchUrl(w http.ResponseWriter, r *http.Request) {
	//Extracting key string from url
	keyStringFromURL := r.URL.Path[1:]
	if len(keyStringFromURL) != 6 {
		w = utils.JSONError(w, errors.New("key length should be of six characters"), http.StatusBadRequest)
		return
	}

	value, isKeyExist, err := h.KeyValueStore.GetKey(context.Background(), keyStringFromURL)
	if err != nil {
		w = utils.JSONError(w, err, http.StatusBadRequest)
		return
	}
	if isKeyExist {
		http.Redirect(w, r, *value, http.StatusPermanentRedirect)
	}

	w = utils.JSONError(w, fmt.Errorf("matching url for key %s not found", keyStringFromURL), http.StatusBadRequest)
}

func (h *Handler) ShortenURL(ctx context.Context, url URL) (*string, error) {
	check := true
	var err error
	var value *string
	var randomKey string
	// Checking the existence of key in db
	for check {
		randomKey = utils.GetRandomString()
		value, check, err = h.KeyValueStore.GetKey(ctx, randomKey)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}
