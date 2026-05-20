package main

import (
	"net/http"
	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/SaujanyaMagarde/go-server/internal/auth"
)

type authHandler func(http.ResponseWriter , *http.Request , database.User)

func (apiCfg *apiConfig) middelwareAuth(handler authHandler) http.HandlerFunc {
	return func (w http.ResponseWriter , r *http.Request){
		api_key , err := auth.GetAPIKey(r.Header);

		if err != nil{
			respondWithError(w,400,"Invalid api key")
			return
		}

		user,err := apiCfg.DB.GetUserByApiKey(r.Context(),api_key)

		if err != nil{
			respondWithError(w,400,"failed to find user")
			return
		}

		handler(w,r,user)
	}
}
