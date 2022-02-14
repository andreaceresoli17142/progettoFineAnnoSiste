package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	"github.com/gorilla/mux"
	"github.com/mingrammer/cfmt"
)

const hostSite = "https://8080-andreaceresoli1-progetto-sqaocv6g7zy.ws-eu31.gitpod.io/"
const sqlServerIp = "172.18.0.1:3306"

var clientId string = ""
var clientSecret string = ""
var redirectUri string = hostSite + "oauth"

func validate(input string) string {
	// remove " ' < > / \ to validate user input
	re := regexp.MustCompile(`[\\\/\<\>\"\']*`)

	return re.ReplaceAllString(input, "")
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

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

func homePage(w http.ResponseWriter, r *http.Request) {
	state := "stringaBella"
	fmt.Println("endpoint hit: home")
	fmt.Fprintf(w, "<a href=\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"> login with paleoId </a> ", clientId, state, redirectUri)
}

func privateArea(w http.ResponseWriter, r *http.Request, email string) {
	fmt.Fprint(w, "Hi ", email, " you are in your private area!")
}

type UsrLoginData struct {
	Salt  int    `db:"salt"`
	PHash string `db:"pHash"`
}

func login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := validate(vars["username"])
	password := validate(vars["password"])

	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	if err != nil {

	}

	defer db.Close()
	var loginData UsrLoginData
	q := fmt.Sprintf("SELECT salt, pHash FROM Users WHERE username = \"%s\";", username)
	err = db.QueryRow(q).Scan(&loginData.Salt, &loginData.PHash)

	if err == sql.ErrNoRows {
		fmt.Fprint(w, "{ 400:\"bad_request\", error:\"username does not exist\" }")
		return
	}

	if err != nil {
		fmt.Fprintf(w, "{ 300:\"backend_error\", error:\"%v\" }", err)
		return
	}

	sum := sha256.Sum256([]byte(password))

	if string(sum[:]) == loginData.PHash {

		fmt.Fprintf(w, "{ 200:\"login_successful\", access_token:\"placeholder\",  refresh_token:\"placeholder\" }")
		return
	}

	fmt.Fprintf(w, "{ 400:\"bad_request\", error:\"wrong username or password\"  }")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/login/{username}/{password}", login)
	// myRouter.HandleFunc("/oauth", paleoIdAuth).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

type UserData struct {
	Id           int
	Username     string
	Email        string
	Date_of_join string
	Salt         int
	PHash        string
}

// func slqTest() {

// 	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	defer db.Close()
// 	var ret UserData
// 	q := "SELECT * FROM Users;"
// 	err = db.QueryRow(q).Scan(&ret.Id, &ret.Username, &ret.Email, &ret.Date_of_join, &ret.Salt, &ret.PHash)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println(ret)
// }

func main() {

	// slqTest()
	// time.Sleep(3 * time.Second)

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

	fmt.Println("GO server started")
	handleRequests()
}
