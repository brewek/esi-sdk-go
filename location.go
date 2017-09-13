package esi

import (
	"errors"
	"fmt"
	"net/http"
)

type Location struct {
	SolarSystemId int `json:"solar_system_id"`
}

func (c *Client) GetCharacterLocation(characterId int) (Location, error) {
	if c.AccessToken == "" {
		return Location{}, errors.New("No access token. Can't get location.")
	}
	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/characters/%d/location/", characterId), nil, true)
	resp := PerformRequest(req)
	var location Location
	UnmarshalResponse(resp, &location)
	return location, nil
}
