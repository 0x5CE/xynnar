package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
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

func filmsGET(w http.ResponseWriter, r *http.Request) {
	val, _ := makeSWAPICall("films", client)

	var films struct{ Results []Film }
	if err := json.Unmarshal(val, &films); err != nil {
		log.Fatal(err)
	}

	sort.Slice(films.Results, func(i, j int) bool {
		return films.Results[i].Release_Date < films.Results[j].Release_Date
	})

	httpJsonResp(films, w)
}

func charactersGET(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/characters/")
	sortParam := r.URL.Query().Get("sort")
	filterParam := r.URL.Query().Get("filter")

	resp, _ := makeSWAPICall("films/"+id, client)

	var film struct{ Characters []string }
	if err := json.Unmarshal(resp, &film); err != nil {
		log.Fatal(err)
	}

	characters, totalHeight, err := fetchCharacters(film.Characters, filterParam)
	if err != nil {

	}
	sortCharacters(sortParam, characters)

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

	httpJsonResp(response, w)
}

func commentsGET(w http.ResponseWriter, r *http.Request) {
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
		var c Comment
		var id int
		err := rows.Scan(&id, &c.Movie_Id, &c.Content, &c.IP, &c.Timestamp)
		if err != nil {
			panic(err)
		}
		comments.Comments = append(comments.Comments, c)
	}
	comments.Count = len(comments.Comments)

	httpJsonResp(comments, w)
}

func commentPOST(w http.ResponseWriter, r *http.Request) {
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
	ip := r.RemoteAddr

	_, err = db.Exec(fmt.Sprintf(`INSERT INTO film_comments (movie_id, comment, commenter_ip, timestamp)
		VALUES (%s, '%s', '%s', (NOW() AT TIME ZONE 'utc'));`, comment.Movie_Id, comment.Comment, ip))
	if err != nil {
		panic(err)
	}

	var response struct {
		Message string
		IP      string
	}
	response.Message = "Comment posted"
	response.IP = ip

	httpJsonResp(response, w)
}
