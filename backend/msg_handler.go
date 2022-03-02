package main

import (// {{{
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)// }}}

type conversation struct {


}

func getConversations(access_token string) ([]string, error) { // {{{

	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/"+dbname)

	if err != nil {
		// fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return false, err
	}

	defer db.Close()
	var loginData UserData
	// q := fmt.Sprintf("SELECT salt, pHash FROM Users WHERE id = (?);", usr_id)
	err = db.QueryRow("SELECT salt, pHash FROM Users WHERE id = (?);", usr_id).Scan(&loginData.Salt, &loginData.PHash)

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
} // }}}


