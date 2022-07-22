package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var db *sql.DB
var client *redis.Client
var err error

func main() {
	db, client, err = dbInit()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	defer client.Close()

	http.HandleFunc("/films", filmsGET)
	http.HandleFunc("/characters/", charactersGET)
	http.HandleFunc("/comments/", commentsGET)
	http.HandleFunc("/comment", commentPOST)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
