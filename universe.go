package esi

import (
	"fmt"
	"net/http"
	"sync"
)

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Planet struct {
	PlanetID int   `json:"planet_id"`
	MoonsIDs []int `json:"moons"`
}

type System struct {
	StarID          int      `json:"star_id"`
	SystemID        int      `json:"system_id"`
	Name            string   `json:"name"`
	Position        Position `json:"position"`
	ConstellationID int      `json:"constellation_id"`
	Planets         []Planet `json:"planets"`
	SecurityClass   string   `json:"security_class"`
	StargatesIDs    []int    `json:"stargates"`
	StationsIDs     []int    `json:"stations"`
}

func (c *Client) GetSystem(systemID int) (system System) {
	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/universe/systems/%d/", systemID), nil, false)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &system)
	return system
}

type Destination struct {
	SystemID   int `json:"system_id"`
	StargateID int `json:"stargate_id"`
}

type Stargate struct {
	StargateID  int         `json:"stargate_id"`
	Name        string      `json:"name"`
	TypeID      int         `json:"type_id"`
	Position    Position    `json:"position"`
	SystemID    int         `json:"system_id"`
	Destination Destination `json:"destination"`
}

func (c *Client) GetStargate(stargateID int) (stargate Stargate) {
	req := c.CreateRequest(http.MethodGet,
		fmt.Sprintf("/universe/stargates/%d/", stargateID), nil, false)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &stargate)
	return stargate
}

func (c *Client) GetConnectingSystems(systemID int, depth int) (map[int]System,
	map[int]Stargate) {
	systems := make(map[int]System)
	stargates := make(map[int]Stargate)
	var waitGroup sync.WaitGroup
	c.getConnectingSystemsMaps(systemID, depth, systems, stargates, &waitGroup)
	waitGroup.Wait()
	return systems, stargates
}

func (c *Client) getConnectingSystemsMaps(systemID int, depth int,
	systems map[int]System, stargates map[int]Stargate,
	waitGroup *sync.WaitGroup) {

	currSystem := c.GetSystem(systemID)
	systems[currSystem.SystemID] = currSystem

	if depth >= 0 {
		waitGroup.Add(len(currSystem.StargatesIDs))

		for _, stargateID := range currSystem.StargatesIDs {
			go func(stargateID int) {
				defer waitGroup.Done()
				currStargate := c.GetStargate(stargateID)
				stargates[currStargate.StargateID] = currStargate
				c.getConnectingSystemsMaps(currStargate.Destination.SystemID, depth-1, systems, stargates, waitGroup)
			}(stargateID)
		}
	}
}

type SystemJumps struct {
	ShipJumps int `json:"ship_jumps"`
	SystemID  int `json:"system_id"`
}

func (c *Client) GetSystemJumps() map[int]SystemJumps {
	systemJumps := make(map[int]SystemJumps)
	systemJumpsResponse := getSystemJumpsResponse(c)
	for _, systemJump := range systemJumpsResponse {
		systemJumps[systemJump.SystemID] = systemJump
	}
	return systemJumps
}

func getSystemJumpsResponse(c *Client) (systemJumps []SystemJumps) {
	req := c.CreateRequest(http.MethodGet, "/universe/system_jumps/", nil, false)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &systemJumps)
	return systemJumps
}

type SystemKills struct {
	NpcKills  int `json:"npc_kills"`
	PodKills  int `json:"pod_kills"`
	ShipKills int `json:"ship_kills"`
	SystemID  int `json:"system_id"`
}

func (c *Client) GetSystemKills() map[int]SystemKills {
	systemKills := make(map[int]SystemKills)
	systemKillsResponse := getSystemKillsResponse(c)

	for _, systemKill := range systemKillsResponse {
		systemKills[systemKill.SystemID] = systemKill
	}
	return systemKills
}

func getSystemKillsResponse(c *Client) (systemKills []SystemKills) {
	req := c.CreateRequest(http.MethodGet, "/universe/system_kills/", nil, false)
	resp := PerformRequest(req)
	UnmarshalResponse(resp, &systemKills)
	return systemKills
}
