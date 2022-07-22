package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func apiResponse(w http.ResponseWriter, r *http.Request) {
	// Set the return Content-Type as JSON like before
	w.Header().Set("Content-Type", "application/json")

	// Change the response depending on the method being requested
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GET method requested"}`))
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "POST method requested"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

var db *sql.DB
var client *redis.Client
var err error

func makeSWAPICall(endpoint string, client *redis.Client) ([]byte, error) {
	var res []byte
	// check if value is cached
	val, err := client.Get(context.Background(), endpoint).Result()
	if err == redis.Nil {
		resp, err := http.Get("https://swapi.dev/api/" + endpoint)
		if err != nil {
			return nil, err
		}
		res, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// cache for future fetches
		err = client.Set(context.Background(), endpoint, string(res), 48*time.Hour).Err()
	} else if err != nil {
		panic(err)
	} else {
		res = []byte(val) // found
	}
	return res, err
}

func main() {
	db, client, err = dbInit()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/films", getFilms)
	http.HandleFunc("/characters/", getCharacters)
	http.HandleFunc("/comments/", getComments)
	http.HandleFunc("/comment", postComment)
	log.Fatal(http.ListenAndServe(":8080", nil))
	return

	//			<test>
	rows, err := db.Query("SELECT * FROM x")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var x int
		if err := rows.Scan(&x); err != nil {
			panic(err)
		}
		fmt.Println("Read from postgres DB1:", x)
	}
	val, err := makeSWAPICall("films", client)
	fmt.Println(val)
	//			</test>
}
