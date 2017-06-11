package security

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/square/go-jose"
)

var privateKey = new(rsa.PrivateKey)
var encrypter jose.Encrypter

// Encrypt encrypts a byte slice using RSA with a private key
func Encrypt(data []byte) ([]byte, error) {
	object, err := encrypter.Encrypt(data)
	if err != nil {
		return nil, err
	}

	encryptedString := object.FullSerialize()
	return []byte(encryptedString), nil
}

// Decrypt decrypts a byte slice using RSA with a private key
func Decrypt(data []byte) ([]byte, error) {
	object, err := jose.ParseEncrypted(string(data))
	if err != nil {
		return nil, err
	}

	decryptedData, err := object.Decrypt(privateKey)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// InitCipherModule initializes the components used for server-side encryption
func InitCipherModule() {
	var err error

	// Generate the private key
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Initialize the encryption mechanism using the new private key
	encrypter, err = jose.NewEncrypter(jose.RSA_OAEP, jose.A128GCM, &privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
}
