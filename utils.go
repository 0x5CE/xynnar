package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func dbInit() (*sql.DB, *redis.Client, error) {
	// postgres
	psqlInfo := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, nil, err
	}
	// create table if it doesn't exist
	val, err := db.Query("select count(*) from pg_tables where tablename='film_comments';")
	var exists int
	if err != nil {
		return nil, nil, err
	}
	defer val.Close()
	if val.Next() {
		err = val.Scan(&exists)
	}
	if err != nil {
		return nil, nil, err
	}
	if exists == 0 {
		log.Println("Creating table")
		_, err := db.Exec(`CREATE TABLE film_comments (
			id SERIAL, movie_id INT, comment VARCHAR, commenter_ip VARCHAR, timestamp TIMESTAMP);`)
		if err != nil {
			return nil, nil, err
		}
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

// Document represents a regular document.
//
// Link to https://google.com
func heightInFeet(heightInCm string) (string, error) {
	cms, err := strconv.ParseFloat(heightInCm, 64)
	if err != nil {
		return "0'0\"", err
	}
	feet := math.Floor(cms / 30.48)
	inches := cms/2.54 - feet*12
	return fmt.Sprintf("%dft %0.2finches", int(feet), inches), nil
}

func swaggerFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.json")
}
