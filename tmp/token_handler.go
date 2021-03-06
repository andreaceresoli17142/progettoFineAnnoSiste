package main

/*
import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getUserIdFromRefreshToken(refresh_token string) (int, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return -1, err
	}

	defer db.Close()
	var ret int
	q := fmt.Sprintf("SELECT userid FROM Token WHERE refreshToken = \"%v\";", refresh_token)
	err = db.QueryRow(q).Scan(&ret)

	if err == sql.ErrNoRows {
		return -1, nil
	}

	if err != nil {
		return -1, err
	}

	return ret, nil
}

//! old function, moved to accessToken_get_usrid
// func getUserIdFromAccessToken(access_token string) (int, error) {

// 	db, err := sql.Open("mysql", databaseString)

// 	if err != nil {
// 		return -1, err
// 	}

// 	defer db.Close()
// 	var ret int
// 	// q := fmt.Sprintf("SELECT userid FROM Token WHERE accessToken = \"%s\";", access_token)
// 	// err = db.QueryRow(q).Scan(&ret)
// 	err = db.QueryRow("SELECT userid FROM Token WHERE accessToken = (?);", access_token).Scan(&ret)
// 	// fmt.Println("sum: " + fmt.Sprintf("%d", ret))

// 	if err == sql.ErrNoRows {
// 		return -1, nil
// 	}

// 	if err != nil {
// 		return -1, err
// 	}

// 	updateLoginDate(ret)
// 	return ret, nil
// }

func accessTokenExists(access_token string) (bool, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return false, err
	}

	defer db.Close()
	var ret string
	q := fmt.Sprintf("SELECT accessToken FROM Token WHERE accessToken = \"%s\";", access_token)
	err = db.QueryRow(q).Scan(&ret)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func refreshTokenExists(refresh_token string) (bool, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return false, err
	}

	defer db.Close()
	var ret string
	q := fmt.Sprintf("SELECT refreshToken FROM Token WHERE refreshToken = \"%s\";", refresh_token)
	err = db.QueryRow(q).Scan(&ret)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func getUserId_Email(userEmail string) (int, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return -1, err
	}

	defer db.Close()
	var ret int
	q := fmt.Sprintf("SELECT id FROM Users WHERE email = \"%s\";", userEmail)
	err = db.QueryRow(q).Scan(&ret)

	if err == sql.ErrNoRows {
		return -1, nil
	}

	if err != nil {
		return -1, err
	}

	return ret, nil
}

func tokenCoupleAlreadyExists(usrId int) (bool, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		// fmt.Print("0")
		return false, err
	}

	defer db.Close()
	var ret int
	q := fmt.Sprintf("SELECT userid FROM Token WHERE userid = \"%v\";", usrId)
	err = db.QueryRow(q).Scan(&ret)
	// fmt.Print(fmt.Sprint(ret))
	if err == sql.ErrNoRows {
		// fmt.Print("1")
		return false, nil
	}

	if err != nil {
		// fmt.Print("2")
		return false, err
	}

	// fmt.Print("3")
	return true, nil
}

func generateTokenCouple(usrId int) (string, int32, string, error) {

	// fmt.Printf("usrId: %v\n", usrId)

	act := ""

	for {

		act = RandomString(64)
		ret, err := accessTokenExists(act)

		if err != nil {
			// fmt.Println("uh oh")
			return "", -1, "", err
		}

		if !ret {
			// fmt.Println("found act")
			break
		}
		// fmt.Printf("searching act: %v",act)
	}

	rft := ""

	for {

		rft = RandomString(64)
		ret, err := refreshTokenExists(rft)
		// ret := true

		if err != nil {
			return "", -1, "", err
		}

		if !ret {
			// fmt.Println("found rt")
			break
		}
		// fmt.Println("searching rt")
	}

	tokenflag, err := tokenCoupleAlreadyExists(usrId)
	if err != nil {
		return "", -1, "", err
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return "", -1, "", err
	}

	defer db.Close()

	q := ""

	expt := int32(time.Now().Unix()) + 3600

	if tokenflag {
		q = fmt.Sprintf("UPDATE Token SET accessToken = \"%s\", expireTime = %d, refreshToken = \"%s\" WHERE userid = %d;", act, expt, rft, usrId)
	} else {
		q = fmt.Sprintf("INSERT INTO Token (userid, accessToken, expireTime, refreshToken) VALUES (%d, \"%s\", %d, \"%s\");", usrId, act, expt, rft)
	}

	_, err = db.Exec(q)

	if err != nil {
		return "", -1, "", err
	}
	return act, expt, rft, nil
}

type actData struct {
	User_id       int    `db:"userid"`
	Access_token  string `db:"accessToken"`
	Refresh_token string `db:"refreshToken"`
	Exp           int32  `db:"expireTime"`
}

func BearerAuthHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}

	return token
}

func accessToken_get_usrid(access_token string) (int, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return -1, err
	}

	defer db.Close()
	var ret actData
	// q := fmt.Sprintf("SELECT userid, accessToken, expireTime, refreshToken FROM Token WHERE accessToken = \"%s\";", access_token)
	err = db.QueryRow("SELECT userid, accessToken, expireTime FROM Token WHERE accessToken = (?);", access_token).Scan(&ret.User_id, &ret.Access_token, &ret.Exp)

	// fmt.Printf( "\nuser id: %s\n", ret.Access_token )

	if err == sql.ErrNoRows {
				// fmt.Print("wattaf3")
		return -1, nil
	}

	if err != nil {
				// fmt.Print("wattaf1")
		return -1, err
	}

	//return true, nil

	if ret.Access_token != "nil" {
		if ret.Access_token == access_token {
			if int32(time.Now().Unix()) < ret.Exp {
				//! update date is broken ( I think )
				// updateLoginDate(ret.User_id)
				// fmt.Print("wattaf2")
				return ret.User_id, nil
			} // else {
			// 	// err = deleteAccessToken(access_token)
			// 	if err != nil {
			// 		return false, nil
			// 	}
			// }
			//fmt.Printf("now: %d, token_exp: %d", int32(time.Now().Unix()), ret.Exp)
		}
	}
	// fmt.Print("wattaf")
	return -1, nil
}

func useRefreshToken(refresh_token string) (string, int32, string, error) {
	// fmt.Printf("ref token: %s", refresh_token)

	usrId, err := getUserIdFromRefreshToken(refresh_token)

	if err != nil {
		return "", -1, "", err
	}

	if usrId == -1 {
		return "", -1, "", err
	}

	return generateTokenCouple(usrId)
}
*/
