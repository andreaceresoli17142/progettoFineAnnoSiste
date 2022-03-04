package main

import (
	// "log"
	"net/smtp"
)

func sendEmail( reciver_email string, subject string, messagge string) error  {

	// email_email := "noreply.64189489@gmail.com"
	// email_password := "xecntudonnptbhfh"
	// email_server := "smtp.gmail.com"
	//:587
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", email_email, email_password, email_server)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{reciver_email}
	// msg := []byte( fmt.Sprintf(`To: %s \r\n Subject: %s \r\n	\r\n %s \r\n` , reciver_email, subject, messagge ))
	msg_string := "To:"+reciver_email+"\r\nSubject: "+subject+"\r\n\r\n"+messagge+"\r\n"
	msg := []byte(msg_string)
	err := smtp.SendMail(email_server +":"+ email_port, auth, reciver_email, to, msg)
	if err != nil {
		return err
	}
	return nil
}
