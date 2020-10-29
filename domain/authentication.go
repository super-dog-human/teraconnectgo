package domain

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
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

publicKey := publicKey()

func publicKey() *rsa.PublicKey {
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
	ctx := request.Context()
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, FailedDatastoreInitialize
	}

	query := datastore.NewQuery("User").Filter("Auth0Sub =", userSubject).Limit(1)
	keys, err := client.GetAll(ctx, query, &users)
	if err != nil {
		return *user, FailedGettingUser
	}

	if len(users) == 0 {
		return *user, UserNotFound
	}

	user = &users[0]
	user.ID = keys[0].Name
	return users[0], nil
}

// UserSubject is return auth0 subject.
func UserSubject(r *http.Request) (string, error) {
	rawHeader := r.Header.Get("Authorization")
	if rawHeader == "" {
		return "", TokenNotFound
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, UnexpectedSigningMethod
		}
		return publicKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}
	return "", InvalidToken
}
