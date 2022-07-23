package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *sql.DB
var client *redis.Client
var err error

// @title Xynnar API
// @version 1.0
// @description This is an API to get Star Wars movies list and comment on them.
// @contact.url https://github.com/0x5CE/xynnar
// @contact.name Muazzam Ali Kazmi
// @contact.email muazzam_ali@live.com
// @license.name Public Domain
// @BasePath /
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, client, err = dbInit()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	defer client.Close()

	http.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json")))
	http.HandleFunc("/docs/swagger.json", swaggerFiles)

	// API endpoints

	http.HandleFunc("/api/films", filmsGET)
	http.HandleFunc("/api/characters/", charactersGET)
	http.HandleFunc("/api/comments/", commentsGET)
	http.HandleFunc("/api/comment", commentPOST)
	http.HandleFunc("/api/comment/", commentPOST)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
