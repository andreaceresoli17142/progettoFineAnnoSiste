package main

import ( //{{{
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/smtp"
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

var validateEmail string = `^[a-zA-Z0-9.!#$%&'*+=?^_{|}~-@]*$`
var validatePass string = `^[a-zA-Z0-9.!#$%&'*+=?^_{|}~-]{8,}$`
var validateUser string = `^[a-zA-Z0-9.!#$%&'*+=?^_{|}~-]{5,}$`

var vPassErr string = "password must be at least 8 characters long"
var vUserErr string = "username must be at least 6 characters long"
var vEmailErr string = "email is not valid"

//! orribile, da cambiare
func validate(input string, regex string) (string, bool) {
	// remove " ' < > / \ to validate user input
	re := regexp.MustCompile(`[\\\/\<\>\"\']*`)

	ok := false

	if regex != "" {
		var err error
		ok, err = regexp.Match(regex, []byte(input))
		// Debugf("ok? %v, err %v", ok, err)
		if err != nil {
			ok = false
		}
	} else {
		ok = true
	}
	// Debugf("ok? %v", ok)

	return re.ReplaceAllString(input, ""), ok
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
	Debugf("Error: %v", msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s"}`, code, msg)
}

func httpSuccess(w http.ResponseWriter, code int, s string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"code": %d, "msg":"%s"}`, code, s)
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

func sendEmail(reciver_email string, subject string, messagge string) error {

	// email_email := "noreply.64189489@gmail.com"
	// email_password := "xecntudonnptbhfh"
	// email_server := "smtp.gmail.com"
	//:587
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", email_email, email_password, email_server)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{reciver_email}
	// msg := []byte( fmt.Sprintf(`To: %s \r\n Subject: %s \r\n	\r\n %s \r\n` , reciver_email, subject, messagge ))
	msg_string := "To:" + reciver_email + "\r\nSubject: " + subject + "\r\n\r\n" + messagge + "\r\n"
	msg := []byte(msg_string)
	err := smtp.SendMail(email_server+":"+email_port, auth, reciver_email, to, msg)
	if err != nil {
		return err
	}
	return nil
}
