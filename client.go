package esi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	protocol  = "https"
	host      = "esi.tech.ccp.is"
	version   = "latest"
	urlBase   = protocol + "://" + host + "/" + version
	userAgent = "esi-sdk-go"
	source    = "tranquility"
	timeout   = 5 * time.Second
)

type Client struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func addDatasource(url string, source string) string {
	if strings.Contains(url, "datasource=") {
		log.Print("datasource already added")
		return url
	} else if strings.Contains(url, "?") {
		url = url + "&"
	} else {
		url = url + "?"
	}
	return url + "datasource=" + source
}

func (c *Client) CreateRequest(method string, path string, body io.Reader,
	useAuthorization bool) *http.Request {

	url := urlBase + path
	url = addDatasource(url, source)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Host", host)
	if useAuthorization {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	}

	return req
}

func PerformRequest(req *http.Request) *http.Response {

	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Response code was %d.", resp.StatusCode)
		if resp.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Body: %s.", bodyBytes)
		}
		for header, value := range resp.Header {
			log.Printf("%s: %s", header, value)
		}
	}

	return resp
}

func UnmarshalResponse(resp *http.Response, v interface{}) {

	if resp.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(bodyBytes, v)
		if err != nil {
			log.Fatal(err)
		}
	}
}
