package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	// "log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	// "github.com/gorilla/mux"
)

// type UsrLoginData struct {
// 	Email string `db:"email"`
// 	Salt  int    `db:"salt"`
// 	PHash string `db:"pHash"`
// }

func backendLogin( username string, password string ) (bool, error) {

	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	if err != nil {
		// fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return false, err
	}

	defer db.Close()
	var loginData UserData
	q := fmt.Sprintf("SELECT salt, pHash FROM Users WHERE username = \"%s\";", username)
	err = db.QueryRow(q).Scan(&loginData.Salt, &loginData.PHash)

	if err == sql.ErrNoRows {
		// fmt.Fprint(w, "{ \"resp_code\":400, error:\"username does not exist\" }")
		return false, nil
	}

	if err != nil {
		// fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return false, nil
	}

	data := []byte(fmt.Sprint(loginData.Salt) + password)

	hash := sha256.Sum256(data)
	sum := fmt.Sprintf("%x", hash[:])

	if sum == loginData.PHash {
		return true, nil
	}
	return false, nil
}

func login(w http.ResponseWriter, r *http.Request) {// {{{
	fmt.Println("endpoint hit: login")

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	username := validate(r.PostForm.Get("username"))
	password := validate(r.PostForm.Get("password"))

	ret, err := backendLogin(username, password)
	
	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if ret == true {

// start transition
	// db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/instanTex_db")

	// if err != nil {
	// 	fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
	// 	return
	// }

	// defer db.Close()
	// var loginData UserData
	// q := fmt.Sprintf("SELECT email, salt, pHash FROM Users WHERE username = \"%s\";", username)
	// err = db.QueryRow(q).Scan(&loginData.Email, &loginData.Salt, &loginData.PHash)

	// if err == sql.ErrNoRows {
	// 	fmt.Fprint(w, "{ \"resp_code\":400, error:\"username does not exist\" }")
	// 	return
	// }

	// if err != nil {
	// 	fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
	// 	return
	// }

	// data := []byte(fmt.Sprint(loginData.Salt) + password)

	// hash := sha256.Sum256(data)
	// sum := fmt.Sprintf("%x", hash[:])

	// if sum == loginData.PHash {
// end transition
		usrId, err := getUserId_Usr(username)

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

	fmt.Fprintln(w, "{ \"resp_code\":33300, error: \"wrong username or password\" }")
	// fmt.Fprintf(w, "{ \"resp_code\":400, error:\"wrong username or password\"  }")
}// }}}

func accessTokenTest(w http.ResponseWriter, r *http.Request) {// {{{
	fmt.Println("endpoint hit: access token test")

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	act := validate(r.PostForm.Get("access_token"))

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
}// }}}

func refreshTokenReq(w http.ResponseWriter, r *http.Request) {// {{{
	fmt.Println("endpoint hit: use refresh token")

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	reft := validate(r.PostForm.Get("refresh_token"))

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
}// }}}

// func changeUserData(w http.ResponseWriter, r *http.Request){

// 	err := r.ParseForm()

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	act := validate(r.PostForm.Get("access_token"))
// 	username := validate(r.PostForm.Get("new_username"))
// 	email := validate(r.PostForm.Get("new_emal"))
// 	new_pw := validate(r.PostForm.Get("new_password"))
// 	old_pw := validate(r.PostForm.Get("password"))

// 	usrId, err := getUserIdFromAccessToken(act)

// 	if err != nil {
// 		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
// 		return
// 	}

// 	if usrId == -1 {
// 		fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid access token\"  }")
// 		return
// 	}
		
// }

func getUserDataReq(w http.ResponseWriter, r *http.Request) {// {{{
	fmt.Println("endpoint hit: get user data")

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	act := validate(r.PostForm.Get("access_token"))

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
		fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid wrong id\"  }")
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":200, username: \"%v\", email: \"%v\", date_of_join: \"%v\" }", username, email, date_of_join)
}// }}}

func signIn(w http.ResponseWriter, r *http.Request) {// {{{
	fmt.Println("endpoint hit: sign in")

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	username := validate(r.PostForm.Get("username"))
	email := validate(r.PostForm.Get("email"))
	password := validate(r.PostForm.Get("password"))

	resp, err := addUser(username, email, password)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return
	}

	if resp == false {
		fmt.Fprint(w, "{ \"resp_code\":400 \"error\":\"username or email already in use\" }")
		return
	}

	fmt.Fprint(w, "{ \"resp_code\":200 \"details\":\"sign in succesfull\" }")
}// }}}
