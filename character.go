package esi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Character struct {
	AncestryID     int       `json:"ancestry_id"`
	Birthday       time.Time `json:"birthday"`
	BloodlineID    int       `json:"bloodline_id"`
	CorporationID  int       `json:"corporation_id"`
	Description    string    `json:"description"`
	Gender         string    `json:"gender"`
	Name           string    `json:"name"`
	RaceID         int       `json:"race_id"`
	AllianceID     int       `json:"alliance_id"`
	SecurityStatus float32   `json:"security_status"`
}

func (character *Character) GetAlliance() Alliance {
	return (&Client{}).GetAlliance(character.AllianceID)
}

func (character *Character) GetCorporation() Corporation {
	return (&Client{}).GetCorporation(character.CorporationID)
}

func (c *Client) GetCharacter(charactersId int) (character Character) {
	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/characters/%d/", charactersId), nil, false)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &character)
	return character
}

func GetCharacters(charactersIDs []int) []Character {

	c := Client{}

	var characters []Character
	var wg sync.WaitGroup
	wg.Add(len(charactersIDs))
	for _, characterID := range charactersIDs {
		go func(characterID int) {
			defer wg.Done()
			characters = append(characters, c.GetCharacter(characterID))
		}(characterID)
	}
	wg.Wait()

	return characters
}

func createArrayString(input []int) io.Reader {

	body := "["

	for _, i := range input {
		if len(body) > 1 {
			body += ","
		}
		body += "\"" + strconv.Itoa(i) + "\""
	}
	body += "]"
	return bytes.NewBufferString(body)
}

func (c *Client) GetAffiliation(charactersIds []int) string {

	req := c.CreateRequest(http.MethodPost, "/characters/affiliation/",
		createArrayString(charactersIds), false)

	resp := PerformRequest(req)

	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	return string(bodyBytes)
}
