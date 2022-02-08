package main

import (
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
var redirectUri string = "https://8080-andreaceresoli1-progetto-sqaocv6g7zy.ws-eu30.gitpod.io/verifyPaleoIdAuth"

func homePage(w http.ResponseWriter, r *http.Request) {
	state := "stringaBella"
	fmt.Println("endpoint hit: home")
	fmt.Fprintf(w, "<a href=\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"> login with paleoId </a> ", clientId, state, redirectUri)
}

func paleoIdAuth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state, code := query.Get("state"), query.Get("code")
	fmt.Println("endpoint hit: paleoId Auth")
	fmt.Printf("state: %v, code: %v", state, code)

	data := url.Values{
		"grant_type": {"authorization_code"},
		"code": {"e65ad004ebf0c3db769d45ce1323e021bf835682c22e941a0ed252c5a3ce69732679679db16f3feaefc41ffafa0a468553695aeab580ec4b977c883c70f8f0eb"},
		"redirect_uri": {"https://sportellohelp.paleo.bg.it/auth"},
		"client_id": "bb3c313088161696bcd66801b9e7abe4",
		"client_secret": "1b989d248ec27a38904d9cea40be843171c55be50b220176d60d36ff5d35abad1ac546934cd453fd53475b9ae565ca0945f5a9eceb5b1f75cadccd692de2039d"

		"name":       {"John Doe"},
		"occupation": {"gardener"},
	}

	resp, err := http.PostForm("https://id.paleo.bg.it/oauth/token", data)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]string

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)

	//fmt.Fprintf(w, "<a href=\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"> login with paleoId </a> ", clientId, state, redirectUri)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/verifyPaleoIdAuth", paleoIdAuth).Methods("GET")

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

	fmt.Printf("type: %T, content: %v", payload["userid"], payload["userid"])

	//clientId = "f6a5cceda6dc53dc6506a94b2f5f4ed1"

	clientId = payload["userid"]

	handleRequests()
}
