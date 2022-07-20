package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
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

func dbInit() (*sql.DB, *redis.Client, error) {
	// postgres
	psqlInfo := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}
	// redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, nil, err
	}
	return db, client, err
}

func makeSWAPICall(endpoint string, client *redis.Client) (string, error) {
	// check if value is cached
	val, err := client.Get(context.Background(), endpoint).Result()
	if err == redis.Nil {
		resp, err := http.Get("https://swapi.dev/api/" + endpoint)
		if err != nil {
			return "makeSWAPICall: http error", err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "makeSWAPICall: reading error", err
		}
		val = string(body)
		// cache for future fetches
		err = client.Set(context.Background(), endpoint, val, 20*time.Second).Err()
	} else if err != nil {
		panic(err)
	}
	return val, err
}

func main() {
	db, client, err := dbInit()
	if err != nil {
		panic(err)
	}

	//			<test>
	rows, err := db.Query("SELECT * FROM x")
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
