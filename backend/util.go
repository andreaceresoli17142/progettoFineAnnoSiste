package main

import (
	"math/rand"
	"regexp"
)


func validate(input string) string {
	// remove " ' < > / \ to validate user input
	re := regexp.MustCompile(`[\\\/\<\>\"\']*`)

	return re.ReplaceAllString(input, "")
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}