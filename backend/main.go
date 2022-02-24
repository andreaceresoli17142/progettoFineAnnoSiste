package main

import (
	// "crypto/sha256"
	"database/sql"
	"strings"

	// "regexp"
	// "math/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

const hostSite = "https://localhost:8080/"
const sqlServerIp = "172.18.0.1:3306"

var clientId string = ""
var clientSecret string = ""
var redirectUri string = hostSite + "oauth"

type OauthResp struct {
	AccessToken string `json:"access_token"`
}

type UsrData struct {
	Email string `json:"email"`
}

func addState(state string) error {

	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	if err != nil {
		return err
	}

	defer db.Close()

	q := fmt.Sprintf("INSERT INTO LoginState (idstring) VALUES (\"%s\");", state)

	_, err = db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func findState(state string) (string, error) {
	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	if err != nil {
		return "false", err
	}

	defer db.Close()

	var str string
	q := fmt.Sprintf("SELECT * FROM LoginState WHERE idstring = \"%s\";", state)
	err = db.QueryRow(q).Scan(&str)

	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return str, nil
}

func remState(state string) error {
	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	if err != nil {
		return err
	}

	defer db.Close()

	q := fmt.Sprintf("DELETE FROM LoginState WHERE idstring = \"%s\";", state)

	_, err = db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func paleoIdAuth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state, code := query.Get("state"), query.Get("code")
	fmt.Println("endpoint hit: paleoId Auth")

	ret, err := findState(state)

	if err != nil {
		fmt.Println(err)
		return
	}

	if ret == "" {
		return
	}

	err = remState(state)

	if err != nil {
		fmt.Println(err)
		return
	}

	payload := strings.NewReader(fmt.Sprintf(`{"grant_type":"%s" , "code":"%s", "redirect_uri":"%s", "client_id":"%s", "client_secret":"%s" }`, "authorization_code", code, redirectUri, clientId, clientSecret))

	req, err := http.NewRequest("POST", "https://id.paleo.bg.it/oauth/token", payload)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	// now handle the response
	var respData OauthResp
	err = json.Unmarshal(body, &respData)

	if err != nil {
		fmt.Println(err)
		return
	}

	url := "https://id.paleo.bg.it/api/v2/user"

	req, err = http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+respData.AccessToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	var resp1Data UsrData
	err = json.Unmarshal(body, &resp1Data)

	if err != nil {
		fmt.Println(err)
		return
	}

	email := resp1Data.Email

	userId, err := getUserId(email)

	if err != nil {
		fmt.Println(err)
		return
	}

	username, email, date_of_join, err := getUserData(userId)

	if err != nil {
		fmt.Println(err)
		return
	}

	if username == "" {
		fmt.Fprint(w, "user not registered yet")
		return
	}

	fmt.Fprintf(w, "private area: \n\tusername: %s \n\temail: %s \n\tdate of join: %s", username, email, date_of_join)
}

func getOauthLink(w http.ResponseWriter, r *http.Request) {
	var state string
	for {
		state = RandomString(15)
		ret, err := findState(state)
		if err != nil {
			fmt.Println(err)
			return
		}
		if ret == "" {
			break
		}
	}
	err := addState(state)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("endpoint hit: home")
	fmt.Fprintf(w, "{resp_code:\"200\"  link:\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"}", clientId, state, redirectUri)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "helo")
	var state string
	for {
		state = RandomString(15)
		ret, err := findState(state)
		if err != nil {
			fmt.Println(err)
			fmt.Println("helo")
			return
		}
		if ret == "" {
			break
		}
	}
	err := addState(state)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("endpoint hit: home")
	fmt.Fprintf(w, "<a href=\"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v\"> login with paleoId </a> ", clientId, state, redirectUri)
}

type UserData struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Date_of_join string `db:"date_of_join"`
	Salt         int    `db:"salt"`
	PHash        string `db:"pHash"`
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/oauth", paleoIdAuth).Methods("GET")
	myRouter.HandleFunc("/login/{username}/{password}", login).Methods("POST")
	myRouter.HandleFunc("/tokentest/{access_token}", accessTokenTest)
	myRouter.HandleFunc("/refreshtoken/{refresh_token}", refreshTokenReq).Methods("POST")
	myRouter.HandleFunc("/getusrdata/{access_token}", getUserDataReq).Methods("GET")
	myRouter.HandleFunc("/signin/{username}/{email}/{password}", signIn).Methods("POST")

	// myRouter.HandleFunc("/oauth", paleoIdAuth).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {

	//fmt.Println(getUserId("pippo.mario@gimelli.com"))

	// for i := 0; i < 100; i++ {
	// 	fmt.Println(RandomString(10))
	// }

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
