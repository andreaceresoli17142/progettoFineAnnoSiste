package main

import ( // {{{
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
) // }}}

var fileDir string

// const hostSite = "http://localhost:8080/"
// const sqlServerIp = "172.18.0.2:3306"

// var clientId string = ""
// var clientSecret string = ""
// var redirectUri string = hostSite + "oauth"

// type OauthResp struct {
// 	AccessToken string `json:"access_token"`
// }

// type UserData struct {
// 	Id           int    `db:"id"`
// 	Username     string `db:"username"`
// 	Email        string `db:"email"`
// 	Date_of_join string `db:"date_of_join"`
// 	Salt         int    `db:"salt"`
// 	PHash        string `db:"pHash"`
// }

func addState(state string) error { // {{{

	db, err := sql.Open("mysql", databaseString)

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
} // }}}

func findState(state string) (string, error) { // {{{
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		log.Println("aaa")
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
} // }}}

func remState(state string) error { // {{{
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return err
	}

	defer db.Close()

	//q := fmt.Sprintf("DELETE FROM LoginState WHERE idstring = \"%s\";", state)

	_, err = db.Exec("DELETE FROM LoginState WHERE idstring = ?", state)
	if err != nil {
		return err
	}
	return nil
} // }}}

func paleoIdAuth(w http.ResponseWriter, r *http.Request) { // {{{
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

	var resp1Data UserData
	err = json.Unmarshal(body, &resp1Data)

	if err != nil {
		fmt.Println(err)
		return
	}

	email := resp1Data.Email

	userId, err := getUserId_Email(email)

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
} // }}}

func getOauthLink(w http.ResponseWriter, r *http.Request) { // {{{
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
} // }}}

func homePage(w http.ResponseWriter, r *http.Request) { // {{{
	// fmt.Fprintf(w, "helo")
	var state string
	for {
		state = RandomString(15)
		ret, err := findState(state)
		if err != nil {
			fmt.Println(err)
			//	fmt.Println("helo")
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
} // }}}

func test(w http.ResponseWriter, r *http.Request) { // {{{
	fmt.Println("endpoint hit: test")

	err := r.ParseForm()

	if err != nil {
		httpError(w, 200, err)
		// fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	act := validate(r.PostForm.Get("act"))
	t, err := getAccessToken_usrid(act)

	if err != nil {
		httpError(w, 200, err)
		// fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}
	fmt.Fprintf(w, "user id: %d", t)
} // }}}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Executing middleware", r.Method)

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Content-Type", "application/json")
			return
		}

		next.ServeHTTP(w, r)
		// log.Println("Executing middleware again")
	})
}

// route endpoints {{{
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	authRouter := myRouter.PathPrefix("/auth").Subrouter()
	myRouter.HandleFunc("/", homePage) //.Schemes("https")
	authRouter.HandleFunc("/oauth", paleoIdAuth).Methods("GET")
	authRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/test", test)
	authRouter.HandleFunc("/userft", refreshTokenReq).Methods("POST")
	myRouter.HandleFunc("/getusrdata", getUserDataReq).Methods("GET")
	myRouter.HandleFunc("/signin", signIn).Methods("POST")
	myRouter.HandleFunc("/change", changeUserData).Methods("POST")
	myRouter.HandleFunc("/getconversations", getConversations).Methods("GET")
	myRouter.HandleFunc("/retrivepw/getotp", send_otp_retrivePassword).Methods("GET")
	myRouter.HandleFunc("/retrivepw/useotp", use_otp_retrivePassword).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(myRouter)))
} // }}}

func main() { // {{{

	// load enviroment variables
	ok := loadEnv()
	// if loading fails exit the program
	if !ok {
		return
	}

	_, fileDir, _, ok = runtime.Caller(1)
	if !ok {
		log.Fatal("error getting file directory")
	}

	Successln("GO server started")
	handleRequests()
} // }}}
