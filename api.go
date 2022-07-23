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
	Episode_Id     int
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

type filmsGETResp struct {
	Count   int
	Results []Film
}

type charactersGETResp struct {
	Count          int
	TotalHeight    string
	TotalHeight_Ft string
	Characters     []Character
}

type commentsGETResp struct {
	Count    int       `json:"count" example:"12"`
	Comments []Comment `json:"comments"`
}

type commentPOSTBody struct {
	Movie_Id int    `json:"movie_id" example:"12"`
	Comment  string `json:"comment" example:"Great movie!"`
}

// @tags		Films
// @Summary     Get films list
// @Description Retreive list of all Star War movies along with title, opening crawl, etc. Retreived in chronological order.
// @Accept      json
// @Produce     json
// @Success     200 {object} filmsGETResp
// @Router      /api/films [get]
func filmsGET(w http.ResponseWriter, r *http.Request, connect Connect) (any, error) {
	val, err := makeSWAPICall("films", connect)
	if err != nil {
		return "Error fetching films", err
	}

	var films filmsGETResp
	err = json.Unmarshal(val, &films)
	if err != nil {
		return "Error in films data", err
	}

	var commentsCount int
	for i, f := range films.Results {
		rows, err := connect.db.Query(`SELECT count(*) FROM film_comments WHERE 
			movie_id=` + strconv.Itoa(f.Episode_Id))
		if err != nil {
			log.Println("filmsGET: Error in retreiving film comments count")
			commentsCount = -1 // not a major error
		} else {
			defer rows.Close()
			if rows.Next() {
				err = rows.Scan(&commentsCount)
			} else {
				commentsCount = 0
			}
		}
		films.Results[i].Total_Comments = commentsCount
	}
	films.Count = len(films.Results)

	sort.Slice(films.Results, func(i, j int) bool {
		return films.Results[i].Release_Date < films.Results[j].Release_Date
	})
	return films, nil
}

// @tags		Characters
// @Summary     Get movie characters
// @Description Retreive characters in a particular movie specified by ID. Also outputs their combined height.
// @Accept      json
// @Produce     json
// @Param       movie_id  path     int true "Movie ID"
// @Param       sort  query     string false "Sort by name, gender, or height. Place '-' before it for descending order"
// @Param       filter  query     string false "Filter by gender"
// @Success     200 {object} charactersGETResp
// @Router      /api/characters/{movie_id} [get]
func charactersGET(w http.ResponseWriter, r *http.Request, connect Connect) (any, error) {
	id := strings.TrimPrefix(r.URL.Path, "/api/characters/")
	sortParam := r.URL.Query().Get("sort")
	filterParam := r.URL.Query().Get("filter")

	resp, err := makeSWAPICall("films/"+id, connect)
	if err != nil {
		return "Error fetching films", err
	}

	var film struct{ Characters []string }
	err = json.Unmarshal(resp, &film)

	if err != nil {
		return "Error in films data", err
	}

	characters, totalHeight, err := fetchCharacters(film.Characters, filterParam, connect)
	if err != nil {
		return "Error fetching characters", err
	}
	sortCharacters(sortParam, characters)

	var response charactersGETResp
	response.Count = len(characters)
	response.TotalHeight = fmt.Sprintf("%d", int(totalHeight))
	response.TotalHeight_Ft, _ = heightInFeet(response.TotalHeight)
	response.Characters = characters

	return response, nil
}

// @tags		Comments
// @Summary     Get all comments about a movie by ID.
// @Description Comments returns either all comments, or comments about a particular movie. For instance, /comments/4
// @Accept      json
// @Produce     json
// @Param       movie_id   path      int  true  "Movie ID"
// @Success     200 {object} commentsGETResp
// @Router      /api/comments/{movie_id} [get]
func commentsGET(w http.ResponseWriter, r *http.Request, connect Connect) (any, error) {
	id := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	log.Println(id)
	var comments commentsGETResp
	var query string
	if id == "" {
		// all
		query = "SELECT * FROM film_comments ORDER BY timestamp DESC"
	} else {
		query = fmt.Sprintf(`SELECT * FROM film_comments 
			WHERE movie_id=%s ORDER BY timestamp DESC `, id)
	}
	rows, err := connect.db.Query(query)
	if err != nil {
		log.Printf("commentsGET: select error")
		return "Internal error", err
	}
	defer rows.Close()
	for rows.Next() {
		var c Comment
		var id int
		err := rows.Scan(&id, &c.Movie_Id, &c.Content, &c.IP, &c.Timestamp)
		if err != nil {
			log.Printf("commentsGET: scan error")
			return "Internal error", err
		}
		comments.Comments = append(comments.Comments, c)
	}
	comments.Count = len(comments.Comments)

	return comments, nil
}

// @tags		Comments
// @Summary     Comment on a movie
// @Description Commment on a movie (public comment)
// @Accept      json
// @Produce     json
// @Param       body body commentPOSTBody true "Movie ID"
// @Success     200 {object} commentsGETResp
// @Router      /api/comment/ [post]
func commentPOST(w http.ResponseWriter, r *http.Request, connect Connect) (any, error) {
	decoder := json.NewDecoder(r.Body)

	var comment commentPOSTBody
	err := decoder.Decode(&comment)
	if err != nil {
		log.Printf("commentPOST: decode error")
		return "Internal error", err
	}
	if len(comment.Comment) > 500 {
		return "Comment must be less than 500 characters", err
	}
	if comment.Movie_Id == 0 || comment.Comment == "" {
		return "movie_id or comment missing", err
	}
	ip := r.RemoteAddr

	_, err = connect.db.Exec(fmt.Sprintf(`INSERT INTO film_comments (movie_id, comment, commenter_ip, timestamp)
		VALUES (%d, '%s', '%s', (NOW() AT TIME ZONE 'utc'));`, comment.Movie_Id, comment.Comment, ip))
	if err != nil {
		log.Printf("commentPOST: insert error")
		return "Internal error", err
	}
	var response struct {
		Message string
		IP      string
	}
	response.Message = "Comment posted"
	response.IP = ip
	return response, nil
}
