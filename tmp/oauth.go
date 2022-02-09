/*
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"net/http"
	"strings"

	// "github.com/gorilla/mux"
	"github.com/mingrammer/cfmt"
)

var clientId string = ""
var clientSecret string = ""
var redirectUri string = hostSite + "oauth"

type OauthResp struct {
	AccessToken string `json:"access_token"`
}

type UsrData struct {
	Email string `json:"email"`
}

func paleoIdAuth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state, code := query.Get("state"), query.Get("code")
	fmt.Println("endpoint hit: paleoId Auth")

	if state == "stringaBella" {
		//fai cose
	}

	payload := strings.NewReader(fmt.Sprintf(`{"grant_type":"%s" , "code":"%s", "redirect_uri":"%s", "client_id":"%s", "client_secret":"%s" }`, "authorization_code", code, redirectUri, clientId, clientSecret))

	req, err := http.NewRequest("POST", "https://id.paleo.bg.it/oauth/token", payload)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		cfmt.Errorf("Error: %s", err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		cfmt.Errorf("Error: %s", err.Error())
		return
	}

	// now handle the response
	var respData OauthResp
	err = json.Unmarshal(body, &respData)

	if err != nil {
		cfmt.Errorf("Error: %s", err.Error())
		return
	}

	url := "https://id.paleo.bg.it/api/v2/user"

	req, err = http.NewRequest("GET", url, nil)

	if err != nil {
		cfmt.Errorf("Error: %s", err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+respData.AccessToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		cfmt.Errorf("Error: %s", err.Error())
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		cfmt.Errorf("Error: %s", err.Error())
		return
	}

	var resp1Data UsrData
	err = json.Unmarshal(body, &resp1Data)

	email := resp1Data.Email

	privateArea(w, r, email)
}
*/