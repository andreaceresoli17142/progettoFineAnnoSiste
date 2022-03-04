package main

import ( // {{{
	"os"

	"github.com/joho/godotenv"
) // }}}

var databaseString string
var hostSite string
var sqlServerIp string
var dbname string

var clientId string
var clientSecret string
var redirectUri string

var email_email string
var email_password string
var email_server string
var email_port string

func loadEnv() bool {
	godotenv.Load(".env")

	hostSite, ok := os.LookupEnv("HOST_SITE")
	if !ok {
		Errorln("missing host site from env variables")
		return false
	}

	//o needs change
	redirectUri = hostSite + "oauth"

	sqlServerIp, ok = os.LookupEnv("SQL_SERVER_IP")
	if !ok {
		Errorln("missing sql server ip from env variables")
		return false
	}

	dbname, ok = os.LookupEnv("DATABASE_NAME")
	if !ok {
		Errorln("missing database name from env variables")
		return false
	}

	databaseString = "root:root@tcp(" + sqlServerIp + ")/" + dbname

	clientId, ok = os.LookupEnv("CLIENT_ID_OAUTH")
	if !ok {
		Errorln("missing client id from env variables")
		return false
	}

	clientSecret, ok = os.LookupEnv("CLIENT_SECRET_OAUTH")
	if !ok {
		Errorln("missing client secret from env variables")
		return false
	}

	email_email, ok = os.LookupEnv("EMAIL_EMAIL")
	if !ok {
		Errorln("missing noreply email from env variables")
		return false
	}

	email_password, ok = os.LookupEnv("EMAIL_PASSWORD")
	if !ok {
		Errorln("missing noreply password from env variables")
		return false
	}

	email_server, ok = os.LookupEnv("EMAIL_SERVER")
	if !ok {
		Errorln("missing noreply server from env variables")
		return false
	}

	email_port, ok = os.LookupEnv("EMAIL_PORT")
	if !ok {
		Errorln("missing noreply port from env variables")
		return false
	}

	Successln("enviroment variables loaded")
	return true
}
