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

func addState(state string) error { // {{{

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`INSERT INTO LoginState VALUES (?, null)`, state)
	if err != nil {
		Debugf("oh hi there")
		return err
	}
	return nil
} // }}}

func findState(state string) (string, error) { // {{{
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return "false", err
	}

	defer db.Close()

	var str string
	err = db.QueryRow(`SELECT idstring FROM LoginState WHERE idstring = ?`, state).Scan(&str)

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

func oauthGetTokenCouple(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: oauth get token couple")

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		httpError(w, 400, "Content Type is not application/json")
		return
	}

	// type ReqData struct {
	// 	State string `json:"state"`
	// }
	// var re ReqData

	// err := httpGetBody(r, &re)

	// if err != nil {
	// 	httpError(w, 500, err)
	// 	return
	// }

	// state, ok := validate(re.State, "")

	err := r.ParseForm()

	if err != nil {
		httpError(w, 500, err)
		return
	}

	vars := mux.Vars(r)
	state, ok := validate(vars["state"], "")

	if !ok {
		httpError(w, 400, "error validating state")
		return
	}

	// get email
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	defer db.Close()

	var str string
	err = db.QueryRow(`SELECT userEmail FROM LoginState WHERE idstring = ?`, state).Scan(&str)

	if err == sql.ErrNoRows {
		httpError(w, 500, err)
		return
	}

	if err != nil {
		httpError(w, 500, err)
		return
	}

	// finished getting email

	usrId, err := getUserId_Email(str)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	act, act_expt, rft, rft_expt, err := generateTokenCouple(usrId)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	err = updateLoginDate(usrId)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	httpSuccessf(w, 200, `"access_token":"%s", "act_expt": %d, "refresh_token":"%s", "rft_expt":%d`, act, act_expt, rft, rft_expt)
}

func signInOauth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: oauth sign in")

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		httpError(w, 400, "Content Type is not application/json")
		return
	}

	type ReqData struct {
		State    string `json:"state"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	var re ReqData

	err := httpGetBody(r, &re)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	// err = r.ParseForm()

	// if err != nil {
	// 	httpError(w, 500, err)
	// 	return
	// }

	username, ok := validate(re.Username, validateUser)
	// Debugln(ok)
	if !ok {
		httpError(w, 400, vUserErr)
		return
	}

	state, ok := validate(re.State, "")
	// Debugln(ok)
	if !ok {
		httpError(w, 400, "error validating state")
		return
	}

	// get email
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	defer db.Close()

	var str string
	err = db.QueryRow(`SELECT userEmail FROM LoginState WHERE idstring = ?`, state).Scan(&str)

	if err == sql.ErrNoRows {
		httpError(w, 500, err)
		return
	}

	if err != nil {
		httpError(w, 500, err)
		return
	}

	// finished getting email

	email, ok := validate(str, "")

	if !ok {
		httpError(w, 400, vEmailErr)
		return
	}

	password, ok := validate(re.Password, validatePass)

	if !ok {
		httpError(w, 400, vPassErr)
		return
	}

	resp, err := addUser(username, email, password)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	if !resp {
		httpError(w, 400, "username or email already in use")
		// fmt.Fprint(w, "{ \"resp_code\":400 \"error\":\"username or email already in use\" }")
		return
	}

	usrId, err := getUserId_Email(email)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	act, act_expt, rft, rft_expt, err := generateTokenCouple(usrId)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	err = updateLoginDate(usrId)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	httpSuccessf(w, 200, `"access_token":"%s", "act_expt": %d,  "refresh_token":"%s", "rft_expt":%d`, act, act_expt, rft, rft_expt)

	// httpSuccess(w, 200, "sign in succesfull")
	// fmt.Fprint(w, "{ \"resp_code\":200 \"details\":\"sign in succesfull\" }")
}

func paleoIdAuth(w http.ResponseWriter, r *http.Request) { // {{{
	query := r.URL.Query()
	state, code := query.Get("state"), query.Get("code")
	fmt.Println("endpoint hit: paleoId Auth")

	ret, err := findState(state)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	if ret == "" {
		return
	}

	// err = remState(state)

	// if err != nil {
	// 	httpError(w, 500, err)
	// 	return
	// }

	payload := strings.NewReader(fmt.Sprintf(`{"grant_type":"%s" , "code":"%s", "redirect_uri":"%s", "client_id":"%s", "client_secret":"%s" }`, "authorization_code", code, redirectUri, clientId, clientSecret))

	req, err := http.NewRequest("POST", "https://id.paleo.bg.it/oauth/token", payload)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		httpError(w, 500, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	// now handle the response
	var respData OauthResp
	err = json.Unmarshal(body, &respData)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	url := "https://id.paleo.bg.it/api/v2/user"

	req, err = http.NewRequest("GET", url, nil)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+respData.AccessToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	var resp1Data UserData
	err = json.Unmarshal(body, &resp1Data)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	email := resp1Data.Email

	// pair state with email

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	defer db.Close()

	_, err = db.Exec("UPDATE LoginState SET userEmail = ? WHERE idstring = ?", email, state)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	// redirect user

	userId, err := getUserId_Email(email)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	username, _, _, err := getUserData(userId)

	if err != nil {
		httpError(w, 500, err)
		return
	}

	if username == "" {
		http.Redirect(w, r, "http://localhost/signUp/oauth.html", http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, "http://localhost/getoauthtk.html", http.StatusMovedPermanently)
} // }}}

func getOauthLink(w http.ResponseWriter, r *http.Request) { // {{{
	fmt.Println("endpoint hit: getlink")
	var state string
	for {
		state = RandomString(15)
		ret, err := findState(state)
		if err != nil {
			httpError(w, 500, err)
			return
		}
		if ret == "" {
			break
		}
	}
	err := addState(state)
	if err != nil {
		httpError(w, 500, err)
		return
	}
	httpSuccessf(w, 200, `"state":"%v", "link":"https://id.paleo.bg.it/oauth/authorize?client_id=%v&response_type=code&state=%v&redirect_uri=%v"`, state, clientId, state, redirectUri)
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

	// err := r.ParseForm()

	// if err != nil {
	w.WriteHeader(http.StatusOK)
	log.SetOutput(w)
	log.Println("ciao")

	// httpError(w, 200, err)
	// fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
	return
	// }

	// act, _ := validate(r.PostForm.Get("act"), "")

	// t, err := getAccessToken_usrid(act)

	// if err != nil {
	// 	httpError(w, 200, err)
	// 	// fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
	// 	return
	// }
	// fmt.Fprintf(w, "user id: %d", t)
} // }}}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Executing middleware", r.Method)
		origin := r.Header["Origin"]
		// fmt.Println("origin", origin)
		if len(origin) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(origin, ","))
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			headers := strings.Join(r.Header["Access-Control-Request-Headers"], ",")
			// fmt.Println("HEADERS", headers)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			return
		}
		next.ServeHTTP(w, r)
		// log.Println("Executing middleware again")
	})
}

// route endpoints {{{
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage) //.Schemes("https")
	myRouter.HandleFunc("/test", test)
	myRouter.HandleFunc("/getusrdata", getUserDataReq).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/signin", signIn).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/change", changeUserData).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/getconversations", getConversations).Methods("GET", "OPTIONS")

	oauthRouter := myRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.HandleFunc("/", paleoIdAuth).Methods("GET", "OPTIONS")
	oauthRouter.HandleFunc("/getlink", getOauthLink).Methods("GET", "OPTIONS")
	oauthRouter.HandleFunc("/signin", signInOauth).Methods("POST", "OPTIONS")
	oauthRouter.HandleFunc("/gettkcoup/{state}", oauthGetTokenCouple).Methods("GET", "OPTIONS")

	authRouter := myRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/userft", refreshTokenReq).Methods("POST", "OPTIONS")

	pwrRouter := myRouter.PathPrefix("/pwr").Subrouter()
	pwrRouter.HandleFunc("/getotp/{email}", send_otp_retrivePassword).Methods("GET", "OPTIONS")
	pwrRouter.HandleFunc("/useotp", use_otp_retrivePassword).Methods("POST", "OPTIONS")

	// headersOk := handlers.AllowedHeaders([]string{"*"})
	// originsOk := handlers.AllowedOrigins([]string{"*"})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PATH", "PUT", "DELETE", "OPTIONS"})
	// log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(myRouter)))

	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(myRouter)))
} // }}}

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	// log.Println("nyaaa")

	ok := loadEnv()
	// if loading fails exit the program
	if !ok {
		return
	}

	_, fileDir, _, ok = runtime.Caller(1)
	if !ok {
		log.Fatal("error getting file directory")
	}
}

func main() { // {{{

	// sendEmail("andrea.ceresoli03", )

	// Debugf("%T", regexp.MustCompile(`[\\\/\<\>\"\']*`))

	// load enviroment variables

	Successln("GO server started")
	handleRequests()
} // }}}
