package esi

import (
	"fmt"
	"net/http"
)

type Contact struct {
	Standing    float32 `json:"standing"`
	ContactType string  `json:"contact_type"`
	ContactID   int     `json:"contact_id"`
	IsWatched   bool    `json:"is_watched"`
}

func (c *Client) GetCharacterContacts(characterId int) (contacts []Contact) {
	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/characters/%d/contacts/", characterId), nil, true)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &contacts)
	return contacts
}
