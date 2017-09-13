package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	esi "lossprevented.pl/esi-sdk-go"
)

// TokenBody is a structure for holding body of /oauth/token request.
// GrantType either GrantTypeRefreshToken or GrantTypeAuthorizationCode
// RefreshToken if GrantType is GrantTypeRefreshToken, this should contain RefreshToken
// Code if GrantType is GrantTypeAuthorizationCode, this should contain the AuthorizationCode
type TokenBody struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Code         string `json:"code,omitempty"`
}

// GrantTypeRefreshToken is a GrantType for requesting client by RefreshToken
// GrantTypeAuthorizationCode is a GrantType for requesting client (and refresh token) by AuthorizationCode
const (
	GrantTypeRefreshToken      = "refresh_token"
	GrantTypeAuthorizationCode = "authorization_code"
)

func prepareRequest(tb TokenBody, eveAuthClientID string, eveAuthSecretKey string) *http.Request {

	b, err := json.Marshal(tb)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost,
		"https://login.eveonline.com/oauth/token",
		bytes.NewBufferString(string(b)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "login.eveonline.com")

	authorizationPlain := []byte(eveAuthClientID + ":" + eveAuthSecretKey)
	authorizationBase64 := base64.StdEncoding.EncodeToString(authorizationPlain)
	req.Header.Add("Authorization", "Basic "+authorizationBase64)

	return req
}

// GetClient returns Client structure when provided with either,
// authorization code and "false" refresh parameter, or
// refresh_token and "true" refresh parameter.
func GetClient(tb TokenBody, eveAuthClientID string, eveAuthSecretKey string) (esi.Client, error) {

	if tb.Code == "" && tb.RefreshToken == "" {
		return esi.Client{}, nil
	}

	req := prepareRequest(tb, eveAuthClientID, eveAuthSecretKey)
	resp := esi.PerformRequest(req)

	var c esi.Client
	esi.UnmarshalResponse(resp, &c)

	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return esi.Client{}, errors.New(string(bodyBytes))
	}

	return c, nil
}

// ownTime needed because this single response contains time in a unusual format
type ownTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05"

func (ot *ownTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ot.Time = time.Time{}
		return
	}
	ot.Time, err = time.Parse(ctLayout, s)
	return
}

// VerifyResponse is a structure for VerifyBearerToken response
type VerifyResponse struct {
	CharacterID          int     `json:"CharacterID"`
	CharacterName        string  `json:"CharacterName"`
	ExpiresOn            ownTime `json:"ExpiresOn"`
	Scopes               string  `json:"Scopes"`
	TokenType            string  `json:"TokenType"`
	CharacterOwnerHash   string  `json:"CharacterOwnerHash"`
	IntellectualProperty string  `json:"IntellectualProperty"`
}

// VerifyBearerToken is used to verify bearer token and get VerifyResponse (containing amongts other: CharacterID)
func VerifyBearerToken(token string) (response VerifyResponse) {

	req, err := http.NewRequest(http.MethodGet, "https://login.eveonline.com/oauth/verify", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Host", "login.eveonline.com")
	resp := esi.PerformRequest(req)
	esi.UnmarshalResponse(resp, &response)
	return response
}
