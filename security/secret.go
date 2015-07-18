package security

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
)

const keyFileName = ".key"

var encryptionKey string

// Initialize the security module by having an AES encryption key
// either by getting an existent one or by generating a new one
func InitCrypto() {
	if _, err := os.Stat(keyFileName); os.IsNotExist(err) {
		encryptionKey = generateSecretKey()
	} else {
		encryptionKey = fetchSecretKey()
	}
}

// Reads the secret key from a configuration file
func fetchSecretKey() string {
	data, err := ioutil.ReadFile(keyFileName)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(data)
}

// Generates a new secret key and also saves it
// in a configuration file
func generateSecretKey() string {
	key := make([]byte, 64)

	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	err = ioutil.WriteFile(keyFileName, key, os.ModeDevice)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(key)
}
