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

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Msg string `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
    }
	type returnMsg struct {
		Id uuid.UUID `json:"id"` 
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
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
		respondWithError(w, http.StatusBadRequest, "Already exists", err)
		return
	} 

	respondWithJSON(w, 201, returnMsg{
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