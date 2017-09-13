package esi

import (
	"fmt"
	"net/http"
)

type Corporation struct {
	CorporationName        string  `json:"corporation_name"`
	Ticker                 string  `json:"ticker"`
	MemberCount            int     `json:"member_count"`
	CEOID                  int     `json:"ceo_id"`
	CorporationDescription string  `json:"corporation_description"`
	TaxRate                float32 `json:"tax_rate"`
	CreatorId              int     `json:"creator_id"`
	URL                    string  `json:"url"`
}

func (c *Client) GetCorporation(corporationID int) (corporation Corporation) {
	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/corporations/%d/", corporationID), nil, false)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &corporation)
	return corporation
}
