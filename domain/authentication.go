package domain

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type AuthErrorCode uint

const (
    TokenNotFound           AuthErrorCode = 1
	UnexpectedSigningMethod AuthErrorCode = 2
    InvalidToken            AuthErrorCode = 3
	FailedGettingUser       AuthErrorCode = 4
	UserNotFound            AuthErrorCode = 5
)

func (e AuthErrorCode) Error() string {
    switch e {
    case TokenNotFound:
        return "token not found in header"
    case UnexpectedSigningMethod:
        return "unexpected token signing"
	case InvalidToken:
		return "invalid token"
	case FailedGettingUser:
		return "failed getting user"
	case UserNotFound:
		return "user not found"
    default:
        return "unknown token error"
    }
}

// PublicKey is return rsa key from pub file.
func PublicKey() *rsa.PublicKey {
	keyData, _ := ioutil.ReadFile("./teraconnect.pub")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	return publicKey
}

// GetCurrentUser is return logged in user
func GetCurrentUser(request *http.Request) (User, error) {
	user := new(User) // for return blank user when error

	userSubject, err := UserSubject(request)
	if err != nil {
		return *user, err
	}

	var users []User
	ctx := appengine.NewContext(request)
	query := datastore.NewQuery("User").Filter("Auth0Sub =", userSubject).Limit(1)
	keys, err := query.GetAll(ctx, &users)
	if err != nil {
		return *user, FailedGettingUser
	}

	if len(users) == 0 {
		return *user, UserNotFound
	}

	user = &users[0]
	user.ID = keys[0].StringID()
	return users[0], nil
}

// UserSubject is return auth0 subject.
func UserSubject(r *http.Request) (string, error) {
	raw_header := r.Header.Get("Authorization")
	if raw_header == "" {
		return "", TokenNotFound
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, UnexpectedSigningMethod
		} else {
			return PublicKey(), nil
		}
	})

	if err != nil {
		return "", UnexpectedSigningMethod
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", InvalidToken
	}
}
