package domain

import (
	"net/http"
	"strings"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// TODO move to infrastructure for at development settings.
const avatarThumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/avatar/{id}.png"

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetAvailableAvatars(request *http.Request) ([]Avatar, error) {
	ctx := appengine.NewContext(request)

	currentUser, err := GetCurrentUser(request)
	fmt.Printf("foobar")
	fmt.Printf("currentUser %v+\n", currentUser)
	// TODO check error

	var avatars []Avatar
	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		id := key.StringID()
		avatars[i].ID = id
		avatars[i].ThumbnailURL = strings.Replace(avatarThumbnailURL, "{id}", id, 1)
	}

	return avatars, nil
}
