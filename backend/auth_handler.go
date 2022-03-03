package main

import ( // {{{
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
) // }}}

// using refresh token {{{
func refreshTokenReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: use refresh token")

	reft := r.Header.Get("refresh-token")
	// err := r.ParseForm()

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// reft := validate(r.PostForm.Get("refresh_token"))

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

// }}}

// generating tokens {{{
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

//}}}

// get userid from email {{{
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

// }}}

// look if token couple already exists {{{
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

// }}}

// get bearer tokens from header {{{
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

//}}}

// get usrid from access tokens {{{
func getAccessToken_usrid(access_token string) (int, error) {

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
				//? update date might be broken
				updateLoginDate(ret.User_id)
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

// }}}

// login (backend) {{{
func backendLogin(usr_id int, password string) (bool, error) {

	db, err := sql.Open("mysql", databaseString)

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
}

// }}}

// login (request) {{{
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: login")

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	username := validate(r.PostForm.Get("username"))
	password := validate(r.PostForm.Get("password"))

	//fmt.Println(username)
	//fmt.Println(password)

	usrId, err := getUserId_Usr(username)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if usrId == -1 {
		fmt.Fprintln(w, "{ \"resp_code\":300, error: \"wrong username or password\" }")
		return
	}

	ret, err := backendLogin(usrId, password)
	// fmt.Println("hjello: "+fmt.Sprint(usrId))

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if ret {

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

	fmt.Fprintln(w, "{ \"resp_code\":300, error: \"wrong username or password\" }")
}

// }}}

// change user data {{{
func changeUserData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: change user data")

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	act := validate(r.PostForm.Get("access_token"))
	new_username := validate(r.PostForm.Get("new_username"))
	new_email := validate(r.PostForm.Get("new_email"))
	new_pw := validate(r.PostForm.Get("new_password"))
	old_pw := validate(r.PostForm.Get("password"))
	// if new_email != "" {
	// 	fmt.Print("bdsjgvfahbfdah,j")
	// }

	// fmt.Println(act)

	usrId, err := getAccessToken_usrid(act)

	// fmt.Println(usrId)
	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if usrId == -1 {
		fmt.Fprintf(w, "{ \"resp_code\":300, error:\"invalid access token\"  }")
		return
	}

	ret, err := backendLogin(usrId, old_pw)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if ret {
		// fmt.Printf("asked to change username:{%s}, password:{%s}, email:{%s}", new_username, new_pw, new_email)

		q := "UPDATE Users SET"

		sumChangesFlag := false

		if new_username != "" {
			usrid_1, err := getUserId_Usr(new_username)
			if err != nil {
				fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
				return
			}
			if usrid_1 != -1 {
				fmt.Fprint(w, "{ \"resp_code\":300, error: \"username already exists\" }")
				return
			}
			q += " username = \"" + new_username + "\","
			sumChangesFlag = true
		}

		if new_email != "" {
			usrid_1, err := getUserId_Email(new_email)
			if err != nil {
				fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
				return
			}
			if usrid_1 != -1 {
				fmt.Fprint(w, "{ \"resp_code\":300, error: \"account using this email already exists\" }")
				return
			}
			q += " email = \"" + new_email + "\","
			sumChangesFlag = true
		}

		if new_pw != "" {
			loginState, err := backendLogin(usrId, new_pw)
			if err != nil {
				fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
				return
			}
			if !loginState {
				fmt.Fprint(w, "{ \"resp_code\":300, error: \"password already in use\" }")
				return
			}

			salt := RandomInt(100000)

			data := []byte(fmt.Sprint(salt) + new_pw)

			hash := sha256.Sum256(data)
			sum := fmt.Sprintf("%x", hash[:])

			q += " salt = " + fmt.Sprintf("%d", salt) + ", pHash = \"" + sum + "\","
			sumChangesFlag = true
		}

		if !sumChangesFlag {
			fmt.Fprintf(w, "{ \"resp_code\":300, error:\"you need to specify a change\"  }")
			return
		}

		q = q[:len(q)-1]

		q += " WHERE id = " + fmt.Sprintf("%d", usrId) + ";"

		db, err := sql.Open("mysql", databaseString)

		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
			return
		}

		defer db.Close()

		_, err = db.Exec(q)

		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
			return
		}

		fmt.Fprint(w, "{ \"resp_code\":200, error: \"data changed successfully\" }")
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":300, error:\"invalid access token\"  }")
}

// }}}

// get user data {{{
func getUserDataReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get use data")

	act := r.Header.Get("access-token")

	// err := r.ParseForm()

	// 	if err != nil {
	// 		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
	// 		return
	// 	}
	// 	act := validate(r.PostForm.Get("access_token"))

	usrId, err := getAccessToken_usrid(act)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if usrId == -1 {
		fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid access token\"  }")
		return
	}

	username, email, date_of_join, err := getUserData(usrId)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if username == "" {
		fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid wrong id\"  }")
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":200, username: \"%v\", email: \"%v\", date_of_join: \"%v\" }", username, email, date_of_join)
}

// }}}

// sign in {{{
func signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: sign in")

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	username := validate(r.PostForm.Get("username"))
	email := validate(r.PostForm.Get("email"))
	password := validate(r.PostForm.Get("password"))

	resp, err := addUser(username, email, password)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if !resp {
		fmt.Fprint(w, "{ \"resp_code\":400 \"error\":\"username or email already in use\" }")
		return
	}

	fmt.Fprint(w, "{ \"resp_code\":200 \"details\":\"sign in succesfull\" }")
}

// }}}
