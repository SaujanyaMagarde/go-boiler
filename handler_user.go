package main

import (
	"encoding/json"
	"net/http"
	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/google/uuid"
	"time"
)

//r.Method (Is it a POST? A GET?)

// r.Body (The actual JSON data they sent, like {"username": "saujanya", "password": "123"})

// r.Header (Security tokens, API keys, browser information)


// Examples of how you use w:

// w.WriteHeader(201) (Tells the client "201 Created - Successfully made your user!")

// w.Write([]byte("User created!")) (Sends the actual text or JSON payload back to the screen)

//you MUST give me a place to put the incoming request (r) and a tool for you to write the response back (w)
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body) //json.NewDecoder is a built-in Go tool designed to stream and read that raw data.

	params := parameters{} //creates an empty box

	err := decoder.Decode(&params) //This is where the actual translation happens

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

	//r.context() -> It’s a digital timer and leash.
	//If the user closes their browser tab, or if their internet disconnects, the Go HTTP server instantly shouts into the walkie-talkie: "CANCEL!"
	//This tells the database: "Hey, if this request dies, stop what you're doing immediately."
	//in node.js we have to hear for close call continusoly then abort the operation
	//Spring introduced Reactive Programming.
	//go is very efficent this to manage.

	if err != nil{
		respondWithError(w,400,"failed to create user")
		return
	}

	respondWithJson(w,200,databaseUserToUser(user));
}

func (apiCfg *apiConfig) handlerGetUserByApikey(w http.ResponseWriter, r *http.Request , user database.User){
	respondWithJson(w,200,databaseUserToUser(user));
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request , user database.User){
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(),database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err != nil{
		respondWithError(w,400,"failed to fetch posts")
		return
	}

	respondWithJson(w,200,databasePostsToPosts(posts))
}