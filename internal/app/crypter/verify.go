package crypter

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Not used
func ComparePasswords(hashedPwd string, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Error().Err(err)
		return false
	}

	return true
}