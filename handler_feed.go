package main

import (
	"encoding/json"
	"net/http"
	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/google/uuid"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request , user database.User){
	type parameters struct{
		Name string `json:"name"`
		Url string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w,400,"Invalid request body no")	
		return
	}

	feed,err := apiCfg.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})

	if err != nil{
		respondWithError(w,400,"failed to create feed")
		return
	}

	respondWithJson(w,200,databaseFeedToFeed(feed));
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request , user database.User){
	feeds,err := apiCfg.DB.GetAllFeeds(r.Context());
	if err != nil{
		respondWithError(w,400,"failed to get feeds");
		return;
	}
	respondWithJson(w, 200, databaseFeedsToFeeds(feeds));
}

func (apiCfg *apiConfig) handlerGetFeedByID(w http.ResponseWriter, r *http.Request , user database.User){
	feeds,err := apiCfg.DB.GetFeedsByUserid(r.Context(),user.ID);
	if err != nil{
		respondWithError(w,400,"failed to get feed");
		return;
	}
	respondWithJson(w, 200, databaseFeedsToFeeds(feeds));
}