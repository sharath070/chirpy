package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sharath070/Chirpy/internal/auth"
	"github.com/sharath070/Chirpy/internal/database"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: %d", cfg.fileSeverHits.Load())
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileSeverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits set to 0"))
}

/***************************/
/********** USERS **********/
/***************************/

type userParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var params userParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Failed to generate the hash password", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Error creating user", err)
		return
	}

	userRes := userResp{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJson(w, http.StatusCreated, userRes)
}

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var params userParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	userResponse := userResp{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJson(w, http.StatusOK, userResponse)
}

/***************************/
/********* CHIRPS **********/
/***************************/

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
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

	userId, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Bad uuid format", err)
		return
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
