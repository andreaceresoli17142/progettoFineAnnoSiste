package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mingrammer/cfmt"
)

var clientId string = ""
var clientSecret string = ""
var redirectUri string = "https://8080-andreaceresoli1-progetto-sqaocv6g7zy.ws-eu30.gitpod.io/oauth"
var authSuccUri string = "https://8080-andreaceresoli1-progetto-sqaocv6g7zy.ws-eu30.gitpod.io/getact"

var stateMap map[string]bool

type OauthResp struct {
	AccessToken string `json:"access_token"`
}

type UsrData struct {
	Email string `json:"email"`
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	state := "stringaBella"
	/*state := RandomString(15)
	for true {

		state := RandomString(15)
		if _, keyExists := stateMap[state]; keyExists {
			//
			break
		}
	}*/
	fmt.Println("endpoint hit: home")
	fmt.Fprintf(w, "<a href=\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"> login with paleoId </a> ", clientId, state, redirectUri)
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

func privateArea(w http.ResponseWriter, r *http.Request, email string) {
	fmt.Fprint(w, "Hi ", email, " you are in your private area!")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/oauth", paleoIdAuth).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("GO server started")

	content, err := ioutil.ReadFile("./oauthTokens.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	clientId = payload["userid"]
	clientSecret = payload["usersecret"]

	handleRequests()
}
