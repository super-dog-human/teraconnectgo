package domain

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type AuthErrorCode uint

const (
	TokenNotFound             AuthErrorCode = 1
	UnexpectedSigningMethod   AuthErrorCode = 2
	InvalidToken              AuthErrorCode = 3
	FailedDatastoreInitialize AuthErrorCode = 4
	FailedGettingUser         AuthErrorCode = 5
	UserNotFound              AuthErrorCode = 6
)

func (e AuthErrorCode) Error() string {
	switch e {
	case TokenNotFound:
		return "token not found in header"
	case UnexpectedSigningMethod:
		return "unexpected token signing"
	case InvalidToken:
		return "invalid token"
	case FailedDatastoreInitialize:
		return "failed datastore initialize"
	case FailedGettingUser:
		return "failed getting user"
	case UserNotFound:
		return "user not found"
	default:
		return "unknown token error"
	}
}

var publicKey *rsa.PublicKey

func loadPublicKey() *rsa.PublicKey {
	if publicKey != nil {
		return publicKey
	}

	keyData, err := ioutil.ReadFile("./public.pem")
	if err != nil {
		log.Printf("%v", err)
		panic("failed to read pem file.")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		log.Printf("%v", err)
		panic("failed to parse pem file.")
	}

	publicKey = key
	return key
}

// ValidTokenClaims returns claims in JWT.
func ValidTokenClaims(r *http.Request) (map[string]interface{}, error) {
	cookie, err := r.Cookie("__Secure-session-token")
	if err != nil {
		return nil, TokenNotFound
	}

	tokenString := cookie.Value
	if tokenString == "" {
		return nil, TokenNotFound
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, UnexpectedSigningMethod
		}
		return loadPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, InvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, InvalidToken
	}

	return claims, nil
}

// ProviderID is return user account provider and subject in JWT.
func ProviderID(r *http.Request) (string, error) {
	claims, err := ValidTokenClaims(r)

	if err != nil {
		return "", err
	}

	providerID := claims["provider"].(string) + "_" + claims["id"].(string)
	return providerID, nil
}
