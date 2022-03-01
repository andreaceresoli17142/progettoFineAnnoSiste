package main

import (
	"fmt"
	// "log"
	"net/smtp"
)

func sendEmail( reciver_email string, subject string, messagge string) error  {

	email := "noreply.64189489@gmail.com"
	// password := "noreply56677890898796g7"
	password := "xecntudonnptbhfh"
	server := "smtp.gmail.com"

	// Choose auth method and set it up
	auth := smtp.PlainAuth("", email, password, server)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{reciver_email}
	// msg := []byte( fmt.Sprintf(`To: %s \r\n Subject: %s \r\n	\r\n %s \r\n` , reciver_email, subject, messagge ))
	msg_string := "To:"+reciver_email+"\r\nSubject: "+subject+"\r\n\r\n"+messagge+"\r\n"
	fmt.Println(msg_string)
	msg := []byte(msg_string)
	err := smtp.SendMail(server+":587", auth, reciver_email, to, msg)
	if err != nil {
		return err
	}
	return nil
}
