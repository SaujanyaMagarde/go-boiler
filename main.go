package main

import (
	"database/sql" // general package for sql
	"fmt"
	"log"
	"net/http" //general package for http 
	"os"
	"time"

	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/go-chi/chi" //chi is a router package
	"github.com/go-chi/cors" //cors is a package to handle cross origin resource sharing
	"github.com/joho/godotenv" //godotenv is a package to load environment variables from .env file
	_ "github.com/lib/pq" //pq is a driver for postgresql database
)

type apiConfig struct{ // struct to hold the database connection
	DB *database.Queries
}

func main(){
	err := godotenv.Load() // Load environment variables from .env file
	if err != nil {
		log.Fatal("error loading .env file")
	}

	portString := os.Getenv("PORT") // Get port from environment variables

	if portString == ""{
		log.Fatal("Port is not found")
	}

	dbURL := os.Getenv("DB_URL") // Get DB_URL from environment variables
	if dbURL == ""{
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres",dbURL) // create db connection
	//this will create a connection pool with go routines 
	//this is different than thread , which consume 1MB/thread memory , 1kb/routines

	if err != nil{
		log.Fatal("error connecting to database")
	}

	queries := database.New(conn)  //passes the connection through queries written by sqlc to initialize them

	//to use this dadabase queries in handler ,we create a sturct which have that connection already 
	//after this we construct all the methods or handler on as a functionof this sturct
	apiCfg := apiConfig{
		DB: queries,
	}

	go startScrapping(apiCfg.DB,10,time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser);
	v1Router.Get("/users", apiCfg.middelwareAuth(apiCfg.handlerGetUserByApikey));
	v1Router.Post("/createfeed",apiCfg.middelwareAuth(apiCfg.handlerCreateFeed));
	v1Router.Get("/getfeeds",apiCfg.middelwareAuth(apiCfg.handlerGetFeeds));
	v1Router.Get("/getfeed",apiCfg.middelwareAuth(apiCfg.handlerGetFeedByID));
	v1Router.Post("/feed_follows",apiCfg.middelwareAuth(apiCfg.handlerCreateFeedFollow));
	v1Router.Get("/feed_follows",apiCfg.middelwareAuth(apiCfg.handlerGetFeedFollows));
	v1Router.Delete("/feed_follows/{feedFollowID}",apiCfg.middelwareAuth(apiCfg.handlerDeleteFeedFollow));
	v1Router.Get("/posts",apiCfg.middelwareAuth(apiCfg.handlerGetPostsForUser));
	router.Mount("/v1", v1Router)

	srv := &http.Server{  //This creates a custom server object.
		Handler: router, //Whenever a request comes in, give it to this router to figure out what to do with it.
		Addr:    ":" + portString, // This tells the server which address (and port) to listen on.
	}

	fmt.Printf("server is running on port: %s\n", portString)
	err = srv.ListenAndServe() //constantly waiting for network requests

	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Port:",portString)
}

//some about handler
//in go handler take two parameter http request , http response .
//in some cases we want to user infor also
//in those case we apply middelware
//in put for those middleware
//input for middelware is a fucntion we wqant to run
//output is brand new, standard HTTP function to the chi router.
//middelware internally run this function pass to middelware and return http function
//to gave response we have to write something on W http.Writerequest