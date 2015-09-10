package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"gost/models"
	"io"
	"time"
)

type Sender struct {
	UserId bson.ObjectId
}

type Device struct {
	Type      string
	OS        string
	Country   string
	State     string
	City      string
	Time      time.Time
	IPAddress string
}

type Credentials struct {
	models.User
	Device
}

// Encrypt data using AES encryption
func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}

	encodedString := base64.StdEncoding.EncodeToString(data)
	encryptedData := make([]byte, aes.BlockSize+len(encodedString))

	iv := encryptedData[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(encryptedData[aes.BlockSize:], []byte(encodedString))

	return encryptedData, nil
}

// Decrypt AES encrypted data
func decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))

	if err != nil {
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(encryptedData, encryptedData)

	decodedBytes, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, err
	}

	return decodedBytes, nil
}
