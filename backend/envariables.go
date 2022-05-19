package main

import ( // {{{

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
) // }}}

var (
	databaseString string
	hostSite       string
	sqlServerIp    string
	dbname         string
	rootcred       string
	updatecred     string
	insertcred     string
	selectcred     string

	clientId     string
	clientSecret string
	redirectUri  string

	email_email    string
	email_password string
	email_server   string
	email_port     string

	rsaPrivateKey rsa.PrivateKey
	rsaPublicKey  rsa.PublicKey
)

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

	databaseString = "@tcp(" + sqlServerIp + ")/" + dbname

	rootcred, ok = os.LookupEnv("CRED_ROOT")
	if !ok {
		Errorln("missing root credentials")
		return false
	}

	updatecred, ok = os.LookupEnv("CRED_UPDATE")
	if !ok {
		Errorln("missing update credentials")
		return false
	}

	insertcred, ok = os.LookupEnv("CRED_INSERT")
	if !ok {
		Errorln("missing insert credentials")
		return false
	}

	selectcred, ok = os.LookupEnv("CRED_SELECT")
	if !ok {
		Errorln("missing select credentials")
		return false
	}

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

	pk, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		// Errorln("missing noreply port from env variables")
		// return false
		kp, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			Errorln("error generating private key")
			return false
		}
		privateKey = kp
		pks := ExportRsaPrivateKeyAsPemStr(kp)
		os.Setenv("PRIVATE_KEY", pks)
	} else {
		rpk, err := ParseRsaPrivateKeyFromPemStr(pk)
		if err != nil {
			Errorln("error parsing private key")
			return false
		}
		privateKey = rpk
	}

	publicKey = &privateKey.PublicKey

	// secretMessage := "0"

	// encryptedMessage, err := rsa_Encrypt(secretMessage, *publicKey)
	// if err != nil {
	// 	Errorln("error parsing private key")
	// 	return false
	// }

	// plainText, err := rsa_Decrypt(encryptedMessage, *privateKey)

	// if plainText != secretMessage || err != nil {
	// 	Errorln("error generating key pais")
	// 	return false
	// }

	// Successln("fetched key pairs")

	Successln("enviroment variables loaded")
	return true
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func rsa_Encrypt(secretMessage string, key rsa.PublicKey) (string, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func rsa_Decrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		return "", err
	}
	// fmt.Println("Plaintext:", string(plaintext))
	return string(plaintext), nil
}

func verifyRsaSignature(publicKey *rsa.PublicKey, keyString string) (string, error) {

	dataSigPair := strings.Split(keyString, ".")

	getTkDigest := sha256.Sum256([]byte(dataSigPair[0]))

	decodedTk, _ := base64.StdEncoding.DecodeString(dataSigPair[1])

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, getTkDigest[:], decodedTk)

	if err != nil {
		return "", err
	}
	return dataSigPair[0], nil
}

func generateRsaSignature(msg []byte, pk *rsa.PrivateKey) (string, error) {

	msgDigest := sha256.Sum256(msg)

	signature, err := rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA256, msgDigest[:])

	if err != nil {
		return "", AppendError("error generating signture", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
