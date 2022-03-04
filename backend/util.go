package main

import ( //{{{
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"
) // }}}

// structs {{{
type OauthResp struct {
	AccessToken string `json:"access_token"`
}

type UserData struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Date_of_join string `db:"date_of_join"`
	Salt         int    `db:"salt"`
	PHash        string `db:"pHash"`
}

type actData struct {
	User_id       int    `db:"userid"`
	Access_token  string `db:"accessToken"`
	Refresh_token string `db:"refreshToken"`
	Act_expt      int    `db:"act_expt"`
	Rft_expt      int    `db:"rft_expt"`
}

type Conversation struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

// }}}

func validate(input string) string {
	// remove " ' < > / \ to validate user input
	re := regexp.MustCompile(`[\\\/\<\>\"\']*`)

	return re.ReplaceAllString(input, "")
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomInt(n int) int {

	rand.Seed(time.Now().UnixNano())

	return rand.Int() % n
}

func httpError(w http.ResponseWriter, code int, msg interface{}) {
	Errorln(msg)
	http.Error(w, fmt.Sprintf(`{"code": %d, "msg":"%s"}`, code, fmt.Sprint(msg)), code)
}

// Access-Control-Allow-Origin: *
// Access-Control-Allow-Methods: POST,GET,PUT,DELETE
// Access-Control-Allow-Headers: Authorization, Lang

func httpSuccess(w http.ResponseWriter, code int, s string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"code": %d, %s}`, code, s)
}

func httpSuccessf(w http.ResponseWriter, code int, s string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	rs := fmt.Sprintf(s, args...)
	fmt.Fprintf(w, `{"code": %d, %s}`, code, rs)
}

func httpGetBody(r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &v)

	if err != nil {
		return err
	}
	return nil
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling Preflight")
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Max-Age", "86400")
	w.WriteHeader(200)
}
