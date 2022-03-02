package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreaceresoli17142/progettoFineAnnoSiste/backend/internal/colorize"
	"github.com/joho/godotenv"
)

var hostSite string // = "http://localhost:8080/"
var sqlServerIp string // = "172.18.0.2:3306"
var dbname string // = "instanTex_db"  

var clientId string // = ""
var clientSecret string // = ""
var redirectUri string //  = hostSite + "oauth"

func loadEnv(){
	godotenv.Load(".env")
	hostSite, ok := os.LookupEnv("HOST_SITE")
	if !ok {
		log.Fatal("missing host site from env variables")
	}

	// needs change
	 redirectUri = hostSite + "oauth"

	sqlServerIp, ok = os.LookupEnv("SQL_SERVER_IP")
	if !ok {
		log.Fatal("missing sql server ip from env variables")
	}

	dbname, ok = os.LookupEnv("DATABASE_NAME")
	if !ok {
		log.Fatal("missing database name from env variables")
	}

	clientId, ok = os.LookupEnv("CLIENT_ID_OAUTH")
	if !ok {
			  log.Fatal("missing client id from env variables")
	}

	clientSecret, ok = os.LookupEnv("CLIENT_SECRET_OAUTH")
	if !ok {
		log.Fatal("missing client secret from env variables")
	}

	// value, ok = os.LookupEnv("")
	// if !ok {
	// 	log.Fatal("from env variables")
	// }
	fmt.Println( cz.Successln("enviroment variables loaded") )
}
