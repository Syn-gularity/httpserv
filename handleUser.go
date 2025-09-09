package main

import (
	"time"
	"github.com/google/uuid"
	"encoding/json"
	"net/http"
	"github.com/Syn-gularity/httpserv/internal/auth"
	"github.com/Syn-gularity/httpserv/internal/database"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request){
    type parameters struct {
		Password string `json:"password"`
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
	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Hash Password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPW,
	})

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

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request){
    type parameters struct {
		Password string `json:"password"`
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

	hash, err := cfg.db.GetPassword(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Couldnt get Hashed PW", err)
		return
	} 

	err = auth.CheckPasswordHash(params.Password, hash)
	if err != nil {
		respondWithError(w, 401, "Wrong Password", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, "Couldn't GetUser", err)
		return
	}

	respondWithJSON(w, 200, returnMsg{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}