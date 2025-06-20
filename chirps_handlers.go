package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sharath070/Chirpy/internal/auth"
	"github.com/sharath070/Chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Body string `json:"body"`
	}

	// check for valid jwt token
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "auth token not found", err)
		return
	}

	userId, err := auth.ValidateJWT(authToken, cfg.jwtSecret)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "unauthorized user", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Error marshalling json", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithErr(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	bodySlice := strings.Fields(params.Body)
	for i := range bodySlice {
		switch bodySlice[i] {
		case "kerfuffle", "sharbert", "fornax", "profane":
			bodySlice[i] = "****"
		}
	}

	chirp, err := cfg.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   strings.Join(bodySlice, " "),
		UserID: userId,
	})
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Error inserting chirp", err)
		return
	}

	res := struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}

	respondWithJson(w, http.StatusOK, res)
}

type chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(context.Background())
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Error getting chirps", err)
		return
	}

	var chirpsRes []chirp
	for _, c := range chirps {
		res := chirp{
			Id:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		}
		chirpsRes = append(chirpsRes, res)
	}

	respondWithJson(w, http.StatusOK, chirpsRes)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n i am in get chipr")
	chirpIdStr := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid uuid format", err)
		return
	}

	c, err := cfg.dbQueries.GetChirp(context.Background(), id)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Couldn't retrive chirp", err)
		return
	}

	chirpRes := chirp{
		Id:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserId:    c.UserID,
	}

	respondWithJson(w, http.StatusOK, chirpRes)
}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid uuid format", err)
		return
	}

	c, err := cfg.dbQueries.DeleteChirp(context.Background(), id)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Couldn't retrive chirp", err)
		return
	}

	chirpRes := chirp{
		Id:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserId:    c.UserID,
	}

	respondWithJson(w, http.StatusOK, chirpRes)
}
