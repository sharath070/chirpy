package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithErr(w http.ResponseWriter, code int, msg string, err error) {
	log.Println(err)

	if code >= 500 {
		log.Printf("Responding with 5XX error: %s\n", msg)
	}

	errResp := struct {
		Error string `json:"error"`
	}{msg}

	respondWithJson(w, code, errResp)
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Println("Error Marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
}
