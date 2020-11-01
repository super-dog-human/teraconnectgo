package domain

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
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

var publicKey = loadPublicKey()

func loadPublicKey() *rsa.PublicKey {
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

	return key
}

// ValidTokenClaims returns claims in JWT.
func ValidTokenClaims(r *http.Request) (map[string]interface{}, error) {
	rawHeader := r.Header.Get("Authorization")
	if rawHeader == "" {
		return nil, TokenNotFound
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, UnexpectedSigningMethod
		}
		return publicKey, nil
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

	log.Printf("claims %v", claims)

	return claims, nil
}

// UserSubject is return user subject in JWT.
func UserSubject(r *http.Request) (string, error) {
	claims, err := ValidTokenClaims(r)

	if err != nil {
		return "", err
	}

	return claims["sub"].(string), nil
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

	query := datastore.NewQuery("User").Filter("AuthSub =", userSubject).Limit(1)
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
