package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/google/uuid"
	"time"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w,400,"Invalid request body no")
		return
	}

	user,err := apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil{
		respondWithError(w,400,"failed to create user")
		return
	}

	respondWithJson(w,200,user)
}