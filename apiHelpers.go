package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func httpJsonResp(response any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	out, _ := json.Marshal(response)
	w.Write(out)
}

func httpJsonError(message string, w http.ResponseWriter) {
	var response struct{ Message string }
	response.Message = message
	out, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	w.Write(out)
}

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

func fetchCharacters(charactersLinks []string, filterParam string) ([]Character, float64, error) {
	var characters []Character
	var totalHeight float64
	for _, x := range charactersLinks {
		characterId := strings.TrimPrefix(x, "https://swapi.dev/api/people/")
		resp, _ := makeSWAPICall("people/"+characterId, client)

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
