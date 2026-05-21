package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("Port is not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres",dbURL)
	if err != nil{
		log.Fatal("error connecting to database")
	}

	queries := database.New(conn)

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

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Printf("server is running on port: %s\n", portString)
	err = srv.ListenAndServe()

	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Port:",portString)
}