package main

import (
	"math/rand"
	"regexp"
	"time"
)

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

// type actData struct {
// 	User_id       int    `db:"userid"`
// 	Access_token  string `db:"accessToken"`
// 	Refresh_token string `db:"refreshToken"`
// 	Exp           int32  `db:"expireTime"`
// }

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
