package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sharath070/Chirpy/internal/auth"
	"github.com/sharath070/Chirpy/internal/database"
)

type userParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResp struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefrestToken string    `json:"refrest_token"`
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

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "error creating jwt token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "error creating refresh token", err)
		return
	}

	refresh, err := cfg.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})

	userResponse := userResp{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefrestToken: refresh.Token,
	}
	respondWithJson(w, http.StatusOK, userResponse)
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "malformed header", err)
		return
	}

	refreshToken, err := cfg.dbQueries.GetUserFromRefreshToken(context.Background(), authToken)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "refresh token not found", err)
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "error creating jwt token", err)
		return
	}

	respondWithJson(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
