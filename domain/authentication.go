package domain

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// PublicKey is return rsa key from pub file.
func PublicKey() *rsa.PublicKey {
	keyData, _ := ioutil.ReadFile("./teraconnect.pub")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	return publicKey
}

// GetCurrentUser is return logged in user
func GetCurrentUser(request *http.Request) (User, error) {
	user := new(User) // for return blank user when error

	userSubject, err := userSubject(request)
	if err != nil {
		return *user, err
	}

	var users []User
	ctx := appengine.NewContext(request)
	query := datastore.NewQuery("User").Filter("Auth0Sub =", userSubject).Limit(1)
	_, err = query.GetAll(ctx, &users)
	if err != nil {
		return *user, err
	}

	if len(users) == 0 {
		return *user, nil
	}

	return users[0], nil
}

func userSubject(r *http.Request) (string, error) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		} else {
			return PublicKey(), nil
		}
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("token is invalid")
	}
}
