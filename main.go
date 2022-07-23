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

type Connect struct {
	db     *sql.DB
	client *redis.Client
}

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

	var db *sql.DB
	var client *redis.Client
	var err error

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, client, err = dbInit()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	defer client.Close()

	connect := Connect{db, client}

	http.Handle("/abc", APIHandler{connect, filmsGET})

	http.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json")))
	http.HandleFunc("/docs/swagger.json", swaggerFiles)

	// API endpoints

	http.Handle("/api/films", APIHandler{connect, filmsGET})
	http.Handle("/api/characters/", APIHandler{connect, charactersGET})
	http.Handle("/api/comments/", APIHandler{connect, commentsGET})
	http.Handle("/api/comment", APIHandler{connect, commentPOST})
	http.Handle("/api/comment/", APIHandler{connect, commentPOST})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
