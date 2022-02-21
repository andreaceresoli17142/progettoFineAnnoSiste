package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	// "log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type UsrLoginData struct {
	Email string `db:"email"`
	Salt  int    `db:"salt"`
	PHash string `db:"pHash"`
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: login")
	vars := mux.Vars(r)
	username := validate(vars["username"])
	password := validate(vars["password"])

	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	defer db.Close()
	var loginData UsrLoginData
	q := fmt.Sprintf("SELECT email, salt, pHash FROM Users WHERE username = \"%s\";", username)
	err = db.QueryRow(q).Scan(&loginData.Email, &loginData.Salt, &loginData.PHash)

	if err == sql.ErrNoRows {
		fmt.Fprint(w, "{ \"resp_code\":400, error:\"username does not exist\" }")
		return
	}

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	data := []byte(fmt.Sprint(loginData.Salt) + password)

	hash := sha256.Sum256(data)
	sum := fmt.Sprintf("%x", hash[:])

	if sum == loginData.PHash {

		usrId, err := getUserId(loginData.Email)

		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		}

		act, expt, rft, err := generateTokenCouple(usrId)

		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
			return
		}

		fmt.Fprintf(w, "{ \"resp_code\":200, access_token:\"%s\", expire_time: %d  refresh_token:\"%s\" }", act, expt, rft)
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":400, error:\"wrong username or password\"  }")
}

func accessTokenTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: access token test")
	vars := mux.Vars(r)
	act := validate(vars["access_token"])

	ret, err := accessTokenValid(act)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	if ret {
		fmt.Fprint(w, "{ \"resp_code\":200 }")
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid access token\"  }")
}

func refreshTokenReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: use refresh token")
	vars := mux.Vars(r)
	reft := validate(vars["refresh_token"])

	act, expt, rft, err := useRefreshToken(reft)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	if act != "" {
		fmt.Fprintf(w, "{ \"resp_code\":200, access_token:\"%s\", expire_time: %d  refresh_token:\"%s\" }", act, expt, rft)
		return
	}

	// fmt.Fprintf(w, "{ \"resp_code\":200, access_token:\"%s\", expire_time: %d  refresh_token:\"%s\" }", act, expt, rft)
	fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid refresh token\" }")
}

func getUserDataReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get user data")
	vars := mux.Vars(r)
	act := validate(vars["access_token"])

	usrId, err := getUserIdFromAccessToken(act)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	if usrId == -1 {
		fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid access token\"  }")
		return
	}

	username, email, date_of_join, err := getUserData(usrId)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	if username == "" {
		fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid access token\"  }")
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":200, username: \"%v\", email: \"%v\", date_of_join: \"%v\" }", username, email, date_of_join)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: sign in")
	vars := mux.Vars(r)
	username := validate(vars["username"])
	email := validate(vars["email"])
	password := validate(vars["password"])

	resp, err := addUser(username, email, password)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	if resp == false {
		fmt.Fprint(w, "{ \"resp_code\":400 \"error\":\"username or email already in use\" }")
		return
	}

	fmt.Fprint(w, "{ \"resp_code\":200 \"details\":\"login succesfull\" }")
}

// func useRefreshToken ()
