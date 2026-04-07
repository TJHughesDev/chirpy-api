package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tjhughesdev/pulse/internal/auth"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401 , "api key sent in wrong format", err )
		return
	}

	if key != cfg.apiKey {
		respondWithError(w, 401 , "issue validating api key", err )
		return
	}


	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode paramters", err)
		return 
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, 204, "event of wrong type", errors.New("event of wrong type"))
		return 
	}

	_, err = cfg.db.UpdateUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "couldn't update user to chirpyred", err)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}

