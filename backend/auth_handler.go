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

	// get refresh token from header
	reft := r.Header.Get("refresh-token")

	// use refresh token and return new token couple
	act, a_expt, rft, r_expt, err := useRefreshToken(reft)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if act != "" {
			  fmt.Fprintf(w, "{ \"resp_code\":200, access_token:\"%s\", act_expt: %d  refresh_token:\"%s\", rft_expt:%d }", act, a_expt, rft, r_expt)
		return
	}

	fmt.Fprintf(w, "{ \"resp_code\":400, error:\"invalid refresh token\" }")
}

func getUserIdFromRefreshToken(refresh_token string) (int, error) {
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return -1, err
	}

	defer db.Close()
	
	var ret actData 
	err = db.QueryRow("SELECT userid, rft_expt FROM Token WHERE refreshToken = (?);", refresh_token).Scan(&ret.User_id, &ret.Rft_expt)

	// return -1 if empty or token is expired
	if err == sql.ErrNoRows && int(time.Now().Unix()) < ret.Rft_expt {
		return -1, nil
	}
	
	// else return error
	if err != nil {
		return -1, err
	}

	return ret.User_id, nil
}

//? function may be usless 
func useRefreshToken(refresh_token string) (string, int, string, int, error) {
	// get userid from refresh token
	usrId, err := getUserIdFromRefreshToken(refresh_token)

	// return error
	if err != nil {
		return "", -1, "", -1, err
	}

	// return -1 if the token does not exists
	if usrId == -1 {
		return "", -1, "", -1, err
	}

	return generateTokenCouple(usrId)
}

// }}}

// generating tokens {{{
func generateTokenCouple(usrId int) (string, int, string, int, error) {
	// generate random string for access token (check if token already exists)
	act := ""

	for {

		act = RandomString(64)
		ret, err := accessTokenExists(act)

		if err != nil {
			return "", -1, "", -1, err
		}

		if !ret {
			break
		}
	}

	// generate random string for access token (check if token already exists)
	rft := ""

	for {

		rft = RandomString(64)
		ret, err := refreshTokenExists(rft)

		if err != nil {
			return "", -1, "", -1, err
		}

		if !ret {
			break
		}
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return "", -1, "", -1, err
	}

	defer db.Close()

	// set access token expire time to 30 min
	act_expt := int(time.Now().Unix()) + 1800
	// set refresh token expire time to 7d
	rft_expt := int(time.Now().Unix()) + 604800
	
	_, err = db.Exec(`
	INSERT INTO Token VALUES ((?), (?), (?), (?), (?)) 
	ON DUPLICATE KEY 
	UPDATE accessToken = (?), act_expt = (?), refreshToken = (?), rft_expt = (?)
	;`, usrId, act, act_expt, rft, rft_expt, act, act_expt, rft, rft_expt )

	if err != nil {
		return "", -1, "", -1, err
	}
	return act, act_expt, rft, rft_expt, nil
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

//! old function no longer in use
// look if token couple already exists {{{
// func tokenCoupleAlreadyExists(usrId int) (bool, error) {

// 	db, err := sql.Open("mysql", databaseString)

// 	if err != nil {
// 		// fmt.Print("0")
// 		return false, err
// 	}

// 	defer db.Close()
// 	var ret int
// 	q := fmt.Sprintf("SELECT userid FROM Token WHERE userid = \"%v\";", usrId)
// 	err = db.QueryRow(q).Scan(&ret)
// 	// fmt.Print(fmt.Sprint(ret))
// 	if err == sql.ErrNoRows {
// 		// fmt.Print("1")
// 		return false, nil
// 	}

// 	if err != nil {
// 		// fmt.Print("2")
// 		return false, err
// 	}

// 	// fmt.Print("3")
// 	return true, nil
// }

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

	return validate(token)
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
	err = db.QueryRow("SELECT userid, accessToken, act_expt FROM Token WHERE accessToken = (?);", access_token).Scan(&ret.User_id, &ret.Access_token, &ret.Act_expt)

	if err == sql.ErrNoRows {
		return -1, nil
	}

	if err != nil {
		return -1, err
	}

	//? i think checking for nil is unnecessary
	// if ret.Access_token != "nil" {

	//check if access token matches 
	if ret.Access_token == access_token {
		// check if token is expired
 		if int(time.Now().Unix()) < ret.Act_expt {
			// update last login date
	 		err = updateLoginDate(ret.User_id)
			if err != nil {
				return -1, err
			}
			return ret.User_id, nil
 		} 
	}
	// }
	return -1, nil
}

// }}}

// login {{{
func backendLogin(usr_id int, password string) (bool, error) {

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return false, err
	}

	defer db.Close()
	var loginData UserData
	err = db.QueryRow("SELECT salt, pHash FROM Users WHERE id = (?);", usr_id).Scan(&loginData.Salt, &loginData.PHash)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, nil
	}

	// compute hash of password and salt
	data := []byte(fmt.Sprint(loginData.Salt) + password)
	hash := sha256.Sum256(data)
	sum := fmt.Sprintf("%x", hash[:])

	if sum == loginData.PHash {
		return true, nil
	}
	return false, nil
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: login")

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	username := validate(r.PostForm.Get("username"))
	password := validate(r.PostForm.Get("password"))

	usrId, err := getUserId_Usr(username)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if usrId == -1 {
		fmt.Fprintln(w, "{ \"resp_code\":400, error: \"wrong username or password\" }")
		return
	}

	ret, err := backendLogin(usrId, password)

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if ret {

		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		}

		act, act_expt, rft, rft_expt, err := generateTokenCouple(usrId)

		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
			return
		}
		
		err = updateLoginDate(usrId)
		
		if err != nil {
			fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
			return
		}

		fmt.Fprintf(w, "{ \"resp_code\":200, act_expt:\"%s\", expire_time: %d  refresh_token:\"%s\", rft_expt:%d }", act, act_expt, rft, rft_expt)
		return
	}

	fmt.Fprintln(w, "{ \"resp_code\":400, error: \"wrong username or password\" }")
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

	act := BearerAuthHeader(r.Header.Get("Authorization"))
	new_username := validate(r.PostForm.Get("new_username"))
	new_email := validate(r.PostForm.Get("new_email"))
	new_pw := validate(r.PostForm.Get("new_password"))
	old_pw := validate(r.PostForm.Get("password"))

	usrId, err := getAccessToken_usrid(act)

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

	act := BearerAuthHeader(r.Header.Get("Authorization"))

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
