package interface

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	//	"github.com/dgrijalva/jwt-go"
	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const thumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/avatar/{id}.png"

// GetAvatars is get lesson avatar.
func GetAvatars(c echo.Context) error {
	// TODO pagination.
	ctx := appengine.NewContext(c.Request())

	//	user := c.Get("user").(*jwt.Token)
	//	claims := user.Claims.(jwt.MapClaims)
	//	name := claims["name"].(string)

	var avatars []domain.Avatar
	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(avatars) == 0 {
		errMessage := "avatars not found"
		log.Warningf(ctx, "%v\n", errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	for i, key := range keys {
		id := key.StringID()
		avatars[i].ID = id
		avatars[i].ThumbnailURL = strings.Replace(thumbnailURL, "{id}", id, 1)
	}

	return c.JSON(http.StatusOK, avatars)
}
