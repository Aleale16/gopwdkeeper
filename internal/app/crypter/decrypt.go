package crypter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"pwdkeeper/internal/app/initconfig"

	"github.com/rs/zerolog/log"
)

type User struct{
	login string
}

func IsAuhtorized(msg string) (login string) {
		user, err := GetAuthenticatedUser(msg)
		if err != nil {
			log.Error().Msgf("Unauhtorized user.login=%v", user.login)
			return ""
		} else {
			log.Info().Msgf("Auhtorized user.login=%v", user.login)
			return user.login
		}
}


func GetAuthenticatedUser(msg string) (user User, err error) {
	var validSign bool	
	log.Debug().Str("Authorization header",msg)
	if msg != "" {
		validSign, user.login = CheckUserAuth(msg)
		if !validSign{
			user.login = ""
			log.Debug().Msgf("Invalid signature!")
			err = errors.New("invalid Signature")
		} else {
			err = nil
		}
	} else {
		user.login = ""
		log.Debug().Msgf("Empty Authorization header!")
		err = errors.New("empty authorization header")
	}	
    //validate the session token in the request,
    //fetch the session state from the session store,
    //and return the authenticated user
    //or an error if the user is not authenticated	
	return user, err
}
func CheckUserAuth(msg string) (validSign bool, val string) {
	var (
		data []byte // декодированное сообщение с подписью
		usefuldata   string // полезное содержимое
		err  error
		sign []byte // HMAC-подпись от полезного содержимого
	)
	validSign = false
	data, err = hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	log.Debug().Str("data", string(data)).Msgf("hex.DecodeString(msg)")
	usefuldata = string(data[sha256.Size:])
	val = usefuldata
	h := hmac.New(sha256.New, initconfig.ServerKey)
	h.Write(data[sha256.Size:])
	sign = h.Sum(nil)
	if hmac.Equal(sign, data[:sha256.Size]) {
		log.Info().Msgf("Подпись подлинная. содержимое:%v", usefuldata)
		validSign = true
	} else {
		log.Warn().Msgf("Подпись неверна. Где-то ошибка! содержимое:%v", usefuldata)
	}
	return validSign, val
}