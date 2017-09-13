package esi

import (
	"fmt"
	"net/http"
	"time"
)

type Alliance struct {
	AllianceName string    `json:"alliance_name"`
	Ticker       string    `json:"ticker"`
	DateFounded  time.Time `json:"date_founded"`
	ExecutorCorp int       `json:"executor_corp"`
}

func (c *Client) GetAlliance(allianceID int) (alliance Alliance) {

	if allianceID == 0 {
		return Alliance{"", "", time.Unix(0, 0), 0}
	}

	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/alliances/%d/", allianceID), nil, false)

	resp := PerformRequest(req)

	UnmarshalResponse(resp, &alliance)

	return alliance
}
