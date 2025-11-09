package masterkey

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	keyLength    = 32
	saltLength   = 16
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
)

type Encrypted struct {
	EncodedCipher string
	EncodedSalt   string
}

type Decrypted struct {
	MasterKey string
}

func Generate() (string, error) {
	key := make([]byte, keyLength)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func Encrypt(masterKey, password string) (*Encrypted, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	derivedKey := deriveKey(password, salt)

	aesGCM, nonce, err := newGCM(derivedKey)
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(masterKey), nil)

	return &Encrypted{
		EncodedCipher: base64.StdEncoding.EncodeToString(ciphertext),
		EncodedSalt:   base64.StdEncoding.EncodeToString(salt),
	}, nil
}

func Decrypt(encryptedMasterKey, password, saltBase64 string) (*Decrypted, error) {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, errors.New("invalid salt encoding")
	}

	derivedKey := deriveKey(password, salt)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherData, err := base64.StdEncoding.DecodeString(encryptedMasterKey)
	if err != nil {
		return nil, errors.New("invalid encrypted data encoding")
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherData) < nonceSize {
		return nil, errors.New("invalid encrypted data length")
	}

	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
	masterKey, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("error", err)
		return nil, errors.New("failed to decrypt master key: authentication failed")
	}

	return &Decrypted{MasterKey: string(masterKey)}, nil
}

func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, keyLength)
}

func newGCM(key []byte) (cipher.AEAD, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	return aesGCM, nonce, nil
}
