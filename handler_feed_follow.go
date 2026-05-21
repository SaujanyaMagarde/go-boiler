package main

import (
	"encoding/json"
	"net/http"
	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request , user database.User){
	type parameters struct{
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w,400,"Invalid request body no")	
		return
	}

	feed_follow,err := apiCfg.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: params.FeedID,
		UserID: user.ID,
	})

	if err != nil{
		respondWithError(w,400,"failed to create feed_follow")
		return
	}

	respondWithJson(w,200,databaseFeedFollowToFeedFollow(feed_follow));
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request , user database.User){
	feed_follows,err := apiCfg.DB.GetFeedFollows(r.Context(),user.ID);
	if err != nil{
		respondWithError(w,400,"failed to get feed_follows")
		return
	}

	respondWithJson(w,200,databaseFeedFollowsToFeedFollows(feed_follows));
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request , user database.User){
	feedFollowIDStr := chi.URLParam(r,"feedFollowID")

	feedFollowID,err := uuid.Parse(feedFollowIDStr)
	if err != nil{
		respondWithError(w,400,"invalid feed follow id")
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(),database.DeleteFeedFollowsParams{
		ID: feedFollowID,
		UserID: user.ID,
	})
	if err != nil{
		respondWithError(w,400,"failed to delete feed follow")
		return
	}
	respondWithJson(w,200,"feed follow deleted")
}