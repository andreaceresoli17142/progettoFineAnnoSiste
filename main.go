package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

var clientId string = ""
var clientSecret string = ""
var redirectUri string = "https://8080-andreaceresoli1-progetto-sqaocv6g7zy.ws-eu30.gitpod.io/oauth"
var authSuccUri string = "https://8080-andreaceresoli1-progetto-sqaocv6g7zy.ws-eu30.gitpod.io/getact"

func homePage(w http.ResponseWriter, r *http.Request) {
	state := "stringaBella"
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
	//fmt.Printf("state: %v, code: %v", state, code)

	body := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {authSuccUri},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	jsonValue, _ := json.Marshal(body)

	fmt.Print(jsonValue)

	resp, err := http.Post("https://id.paleo.bg.it/oauth/token", "application/json", bytes.NewBuffer(jsonValue))
	//resp, err := http.PostForm("https://id.paleo.bg.it/oauth/token", body)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]string

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)

	//fmt.Fprintf(w, "<a href=\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"> login with paleoId </a> ", clientId, state, redirectUri)
}

func getAccesToken(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Println(query)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/oauth", paleoIdAuth).Methods("GET")
	myRouter.HandleFunc("/getact", getAccesToken) //.Methods("GET")

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
