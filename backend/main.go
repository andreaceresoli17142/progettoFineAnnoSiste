package main

import ( // {{{

	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
) // }}}

var fileDir string

func homePage(w http.ResponseWriter, r *http.Request) { // {{{

	httpSuccess(&w, 200, "hey, this is a homepage")
}

// }}}

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

func testData(w http.ResponseWriter, r *http.Request) {
	socketSendNotification(1, `{"ConvId":"P1","Id": 3, "UserId": 1, "Username": "pima", "Text": "you have been hax", "Time": "2022-05-13 20:36:38"}`)

	// act := BearerAuthHeader(r.Header.Get("Authorization"))

	// uid, err := getAccessToken_usrid(act)

	// if err != nil {
	// 	httpError(&w, 500, "backend error: "+err.Error())
	// 	return
	// }

	// if uid == -1 {
	// 	httpError(&w, 400, "user does not exist")
	// 	return
	// }

	// httpSuccess(&w, 200, fmt.Sprint(uid))
}

// route endpoints {{{
func handleRequests() {
	rootRouter := mux.NewRouter().StrictSlash(true)

	rootRouter.HandleFunc("/", homePage) //.Schemes("https")
	rootRouter.HandleFunc("/getusrdata", getUserDataReq).Methods("GET", "OPTIONS")
	rootRouter.HandleFunc("/signin", signIn).Methods("POST", "OPTIONS")
	rootRouter.HandleFunc("/change", changeUserData).Methods("POST", "OPTIONS")
	rootRouter.HandleFunc("/websock", initSocket).Methods("GET", "OPTIONS")
	rootRouter.HandleFunc("/test", testData).Methods("GET", "OPTIONS")
	rootRouter.HandleFunc("/user/getdata", getUserDataEp).Methods("GET", "OPTIONS")

	//TODO: https://github.com/dghubble/gologin
	oauthRouter := rootRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.HandleFunc("/", paleoIdAuth).Methods("GET", "OPTIONS")
	oauthRouter.HandleFunc("/getlink", getOauthLink).Methods("GET", "OPTIONS")
	oauthRouter.HandleFunc("/signin", signInOauth).Methods("POST", "OPTIONS")
	oauthRouter.HandleFunc("/gettkcoup/{state}", oauthGetTokenCouple).Methods("GET", "OPTIONS")

	authRouter := rootRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/userft", refreshTokenReq).Methods("POST", "OPTIONS")

	pwrRouter := rootRouter.PathPrefix("/pwr").Subrouter()
	pwrRouter.HandleFunc("/getotp/{email}", send_otp_retrivePassword).Methods("GET", "OPTIONS")
	pwrRouter.HandleFunc("/useotp", use_otp_retrivePassword).Methods("POST", "OPTIONS")

	freqRouter := rootRouter.PathPrefix("/freq").Subrouter()
	freqRouter.HandleFunc("/makereq", makeFriendRequest).Methods("POST", "OPTIONS")
	freqRouter.HandleFunc("/accreq", acceptFriendRequest).Methods("POST", "OPTIONS")
	freqRouter.HandleFunc("/getreq", getFriendRequest).Methods("GET", "OPTIONS")

	rootRouter.HandleFunc("/getconv", getConversations).Methods("GET", "OPTIONS")
	rootRouter.HandleFunc("/getsingleconv", getSingleConversation).Methods("GET", "OPTIONS")

	groupRouter := rootRouter.PathPrefix("/group").Subrouter()
	groupRouter.HandleFunc("/adduser", addToGroup).Methods("POST", "OPTIONS")
	groupRouter.HandleFunc("/change", changeGroupData).Methods("POST", "OPTIONS")
	groupRouter.HandleFunc("/getdata", getGroupData).Methods("GET", "OPTIONS")
	groupRouter.HandleFunc("/create", createGroup).Methods("POST", "OPTIONS")
	groupRouter.HandleFunc("/quit", quitGroup).Methods("POST", "OPTIONS")
	groupRouter.HandleFunc("/admin/kick", adminKickUser).Methods("POST", "OPTIONS")
	groupRouter.HandleFunc("/admin/manage", adminManage).Methods("POST", "OPTIONS")

	msgRouter := rootRouter.PathPrefix("/msg").Subrouter()
	msgRouter.HandleFunc("/send", sendMessage).Methods("POST", "OPTIONS")
	msgRouter.HandleFunc("/read", getMessages).Methods("GET", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(rootRouter)))
} // }}}

func init() {

	ok := loadEnv()

	// if enviroment variables loading fails exit the program
	if !ok {
		return
	}

	_, fileDir, _, ok = runtime.Caller(1)
	if !ok {
		log.Fatal("error getting file directory")
	}

	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey = &privateKey.PublicKey
}

func main() {
	Successln("GO server started")

	handleRequests()
}
