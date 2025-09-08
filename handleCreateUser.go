package main

import (
	"time"
	"github.com/google/uuid"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Email string `json:"email"`
    }
	type returnMsg struct {
		Id uuid.UUID `json:"id"` 
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Already exists", err)
		return
	} 

	//respondWithJSON(w, 201,user)
	respondWithJSON(w, 201, returnMsg{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}