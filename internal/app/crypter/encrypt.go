package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"pwdkeeper/internal/app/initconfig"

	"golang.org/x/crypto/pbkdf2"
)

// GenAuthToken uses predefined ServerKey, that compiles on both client+server. Client generate hmac signed login, Server checks if sign valid (made on ServerKey)
func GenAuthToken(login string) (authToken string) {
	//sign user with HMAC, using SHA256
	h := hmac.New(sha256.New, initconfig.ServerKey)
	h.Write([]byte(login))
	dst := h.Sum(nil)
	authToken = string(dst) + string(login)
	return hex.EncodeToString([]byte(authToken))
}

/*
At registration time, a uniform random symmetric key (key1) is generated from a CSPRNG using client-side scripting running in the user's web browser.
The user chooses a good password ( preferably dicewire or Bip39),
and another symmetric key (key2) is derived from this password using a password based key derivation function (such as PBKDF2, Argon2, etc).
Then, key1 is encrypted using key2 in-browser, then encrypted key1 is stored on the server.

When the user logs in, encrypted key1 is downloaded to the user's web browser.
The user provides the password, key2 is derived from the password, and encrypted key1 is decrypted in-browser.
key1 is used to encrypt all the user's secrets in-browser, so that only the encrypted secrets are uploaded to the server.
*/
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
//KEK
func Key2build(password string) []byte {
	dk := pbkdf2.Key([]byte(password), initconfig.Salt, 4096, 32, sha1.New)
	return dk
}
//FEK
func Key1build() []byte {
	// будем использовать AES256, создав ключ длиной 32 байта
	key, err := generateRandom(2 * aes.BlockSize) // ключ шифрования
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return key
}

func EncryptKey1(key1, key2 []byte) (key1enc []byte) {
	aesblock, err := aes.NewCipher(key2)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	dst := aesgcm.Seal(nil, nonce, key1, nil) // зашифровываем
	noncedst := append(nonce[:], dst[:]...)
	//fmt.Printf("noncedst %v\n", noncedst)
	return noncedst
}

func DecryptKey1(noncekey1enc, key2 []byte) (key1decrypted []byte) {
	// выделяем вектор инициализации
	nonce := noncekey1enc[:12]
	key1enc := noncekey1enc[12:]

	aesblock, err := aes.NewCipher(key2)
	if err != nil {
		fmt.Printf("error NewCipher: %v\n", err)
		return nil
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error NewGCM: %v\n", err)
		return nil
	}

	key1, err := aesgcm.Open(nil, nonce, key1enc, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error Open: %v\n", err)
		return nil
	}
	return key1
}


func EncryptData (somedata string, key1 []byte)(somedataenc []byte){
	aesblock, err := aes.NewCipher(key1)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	dst := aesgcm.Seal(nil, nonce, []byte(somedata), nil) // зашифровываем
	noncedst := append(nonce[:], dst[:]...)
	return noncedst

}

func DecryptData(noncedataenc, key1 []byte) (somedatadecrypted []byte) {
	// выделяем вектор инициализации
	nonce := noncedataenc[:12]
	dataenc := noncedataenc[12:]

	aesblock, err := aes.NewCipher(key1)
	if err != nil {
		fmt.Printf("error NewCipher: %v\n", err)
		return nil
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error NewGCM: %v\n", err)
		return nil
	}

	somedata, err := aesgcm.Open(nil, nonce, dataenc, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error Open: %v\n", err)
		return nil
	}

	return somedata
}