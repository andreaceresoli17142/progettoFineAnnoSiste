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

	data := []byte(fmt.Sprint(loginData.Salt) + password)

	hash := sha256.Sum256(data)
	sum := fmt.Sprintf("%x", hash[:])

	if sum == loginData.PHash {

		fmt.Fprintf(w, "{ 200:\"login_successful\", access_token:\"placeholder\",  refresh_token:\"placeholder\" }")
		return
	}

	fmt.Fprintf(w, "{ 400:\"bad_request\", error:\"wrong username or password\"  }")
}
