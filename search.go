package esi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type SearchResponse struct {
	Characters []int `json:"character"`
}

// GetCharacterID returns the id for given characterName.
// Expects one result, throws error otherwise.
func (c *Client) GetCharacterID(characterName string) (int, error) {
	searchTerm := strings.Replace(characterName, " ", "%20", -1)
	path := "/search/?categories=character&language=en-us&strict=true&search=" + searchTerm
	req := c.CreateRequest(http.MethodGet, path, nil, false)
	resp := PerformRequest(req)

	var searchResponse SearchResponse
	UnmarshalResponse(resp, &searchResponse)
	if len(searchResponse.Characters) != 1 {
		return -1, errors.New(fmt.Sprintf(
			"wrong amount of characters found, expected 1, got %d",
			len(searchResponse.Characters)))
	}
	return searchResponse.Characters[0], nil
}

func (c *Client) GetCharacterIDs(characterNames []string) []int {

	var characterIDs []int
	var wg sync.WaitGroup
	wg.Add(len(characterNames))

	for _, characterName := range characterNames {
		go func(characterName string) {
			defer wg.Done()
			id, err := c.GetCharacterID(characterName)
			if err != nil {
				log.Printf("%s for characterName %s", err, characterName)
			}
			characterIDs = append(characterIDs, id)
		}(characterName)
	}

	wg.Wait()

	return characterIDs
}
