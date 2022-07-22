package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Film struct {
	Title          string
	Opening_Crawl  string
	Release_Date   string
	Total_Comments int
}

type Character struct {
	Name       string
	Height     string
	Height_Ft  string
	Gender     string
	Birth_Year string
	Hair_Color string
}

type Comment struct {
	Movie_Id  string
	Content   string
	IP        string
	Timestamp time.Time
}

func getFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	val, _ := makeSWAPICall("films", client)

	var films struct{ Results []Film }
	if err := json.Unmarshal(val, &films); err != nil {
		log.Fatal(err)
	}

	sort.Slice(films.Results, func(i, j int) bool {
		return films.Results[i].Release_Date < films.Results[j].Release_Date
	})

	out, _ := json.Marshal(films)
	w.Write(out)
}

func getCharacters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/characters/")

	val, _ := makeSWAPICall("films/"+id, client)

	sortParam := r.URL.Query().Get("sort")
	filterParam := r.URL.Query().Get("filter")

	var film struct{ Characters []string }

	if err := json.Unmarshal(val, &film); err != nil {
		log.Fatal(err)
	}

	var characters []Character
	var totalHeight float64
	for _, x := range film.Characters {
		characterId := strings.TrimPrefix(x, "https://swapi.dev/api/people/")
		val, _ = makeSWAPICall("people/"+characterId, client)
		var character Character
		if err := json.Unmarshal(val, &character); err != nil {
			log.Fatal(err)
		}
		character.Height_Ft, _ = heightInFeet(character.Height)
		if filterParam == "" || character.Gender == filterParam {
			characters = append(characters, character)

			height, _ := strconv.ParseFloat(character.Height, 64)
			totalHeight += height
		}
	}

	asc := true
	if sortParam != "" && sortParam[0] == '-' {
		asc = false
		sortParam = strings.TrimPrefix(sortParam, "-")
	}

	if sortParam == "name" {
		sort.Slice(characters, func(i, j int) bool {
			return asc == (characters[i].Name < characters[j].Name)
		})
	} else if sortParam == "height" {
		sort.Slice(characters, func(i, j int) bool {
			return asc == (characters[i].Height < characters[j].Height)
		})
	} else if sortParam == "gender" {
		sort.Slice(characters, func(i, j int) bool {
			return asc == (characters[i].Gender < characters[j].Gender)
		})
	}

	var response struct {
		Count          int
		TotalHeight    string
		TotalHeight_Ft string
		Characters     []Character
	}
	response.Count = len(characters)
	response.TotalHeight = fmt.Sprintf("%d", int(totalHeight))
	response.TotalHeight_Ft, _ = heightInFeet(response.TotalHeight)
	response.Characters = characters

	out, _ := json.Marshal(response)
	w.Write(out)
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var comments struct {
		Count    int
		Comments []Comment
	}
	rows, err := db.Query("SELECT * FROM film_comments ORDER BY timestamp DESC")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		var id int
		if err := rows.Scan(&id, &comment.Movie_Id,
			&comment.Content, &comment.IP, &comment.Timestamp); err != nil {
			panic(err)
		}
		comments.Comments = append(comments.Comments, comment)
	}

	out, _ := json.Marshal(comments)
	w.Write(out)
}

func postComment(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	fmt.Println(ip)
	decoder := json.NewDecoder(r.Body)
	var comment struct {
		Comment  string
		Movie_Id string
	}
	err := decoder.Decode(&comment)
	if err != nil {
		panic(err)
	}
	if len(comment.Comment) > 500 {
		panic(0)
	}
	if comment.Movie_Id == "" || comment.Comment == "" {
		panic(0)
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO film_comments (movie_id, comment, commenter_ip, timestamp) VALUES (%s, '%s', '%s', (NOW() AT TIME ZONE 'utc'));", comment.Movie_Id, comment.Comment, ip))
	if err != nil {
		fmt.Println("abc")
		panic(err)
	}
	log.Println(comment)
}
