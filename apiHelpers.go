package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type errorResp struct {
	Message   string `json:"message" example:"Error fetching film"`
	ErrorCode int    `json:"count" example:"500"`
}

type APIHandler struct {
	connect  Connect
	endpoint func(http.ResponseWriter, *http.Request, Connect) (any, error)
}

func (handle APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := handle.endpoint(w, r, handle.connect)
	if err != nil { // error response
		var response errorResp
		response.Message = resp.(string)
		response.ErrorCode = 500
		out, _ := json.Marshal(response)

		log.Println(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.ErrorCode)
		w.Write(out)
	} else { // normal response
		w.Header().Set("Content-Type", "application/json")
		out, _ := json.Marshal(resp)
		w.Write(out)
	}
}

func makeSWAPICall(endpoint string, c Connect) ([]byte, error) {
	var res []byte
	// check if value is cached
	val, err := c.client.Get(context.Background(), endpoint).Result()
	if err == redis.Nil {
		log.Println("makeSWAPICall: cache not found, fetching", endpoint)
		resp, err := http.Get("https://swapi.dev/api/" + endpoint)
		if err != nil {
			return nil, err
		}
		res, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// cache for future fetches
		err = c.client.Set(context.Background(), endpoint, string(res), 48*time.Hour).Err()
	} else if err != nil {
		return res, nil
	} else {
		res = []byte(val) // found
	}
	return res, nil
}

func sortCharacters(sortParam string, characters []Character) {
	asc := true
	if strings.HasPrefix(sortParam, "-") {
		asc = false
		sortParam = strings.TrimPrefix(sortParam, "-")
	}

	sort.Slice(characters, func(i, j int) bool {
		var order bool
		switch sortParam {
		case "name":
			order = characters[i].Name < characters[j].Name
		case "height":
			order = characters[i].Height < characters[j].Height
		case "gender":
			order = characters[i].Gender < characters[j].Gender
		}
		return asc == order
	})
}

func fetchCharacters(charactersLinks []string, filterParam string, c Connect) ([]Character, float64, error) {
	var characters []Character
	var totalHeight float64
	for _, url := range charactersLinks {
		characterId := strings.TrimPrefix(url, "https://swapi.dev/api/people/")
		resp, _ := makeSWAPICall("people/"+characterId, c)

		var character Character
		if err := json.Unmarshal(resp, &character); err != nil {
			return characters, totalHeight, err
		}
		character.Height_Ft, _ = heightInFeet(character.Height)
		if filterParam == "" || character.Gender == filterParam {
			characters = append(characters, character)

			height, _ := strconv.ParseFloat(character.Height, 64)
			totalHeight += height
		}
	}
	return characters, totalHeight, nil
}
