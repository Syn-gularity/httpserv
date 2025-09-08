package main

import(
	"time"
	"github.com/google/uuid"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"github.com/Syn-gularity/httpserv/internal/database"
)

type Chirp struct {
	Id uuid.UUID `json:"id"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.NullUUID `json:"user_id"`
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request){
	msg, err := cfg.db.GetMessages(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't run GetMessages Query: ", err)
		return
	} 

	ret := make([]Chirp,len(msg))
	for i:=0; i<len(msg);i++{
		ret[i].Id = msg[i].ID
		ret[i].CreatedAt = msg[i].CreatedAt
		ret[i].UpdatedAt = msg[i].UpdatedAt
		ret[i].Body = msg[i].Body
		ret[i].UserID = msg[i].UserID
	}

	respondWithJSON(w, 200, ret)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request){
	UserID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	msg, err := cfg.db.GetMessage(r.Context(), UserID)
	if err != nil{
		respondWithError(w, 404, "Couldn't run GetMessages Query: ", err)
		return
	} 

	respondWithJSON(w, 200, Chirp{
		Id: msg.ID,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		Body: msg.Body,
		UserID: msg.UserID,
	})
}

func (cfg *apiConfig) handlePostChirps(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Msg string `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Msg) > 140{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	} 
	cleanedMsg := badWordDisplacer(params.Msg)

	msg, err := cfg.db.CreateMessage(r.Context(), database.CreateMessageParams{
		Body: cleanedMsg, 
		UserID: params.UserID,
	})
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	} 

	respondWithJSON(w, 201, Chirp{
		Id: msg.ID,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		Body: msg.Body,
		UserID: msg.UserID,
	})
}

func badWordDisplacer(text string) string{
	ret := ""
	const displacer = "****"
	var badWords = [3]string{"kerfuffle", "sharbert", "fornax"}
	lowered := strings.ToLower(text)
	splitText := strings.Split(lowered," ")
	splitTextOriginal := strings.Split(text," ")
	for idx, word := range splitText{
		bad := false
		for _, badWord := range badWords{
			if word == badWord{
				bad = true
			}
		}
		if bad {
			apnd := fmt.Sprintf(" %s", displacer)
			ret += apnd
		} else {
			apnd := fmt.Sprintf(" %s", splitTextOriginal[idx])
			ret += apnd
		}
	}
	ret = strings.TrimSpace(ret)
	return ret
}