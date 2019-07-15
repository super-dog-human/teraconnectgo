package handler

import (
	"net/http"
	"strings"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const graphicThumbnailURL = "https://storage.googleapis.com/teraconn_thumbnail/graphic/{id}.{fileType}"

func getGraphics(c echo.Context) error {
	// TODO pagination.
	ctx := appengine.NewContext(c.Request())

	var graphics []domain.Graphic
	query := datastore.NewQuery("Graphic").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &graphics)
	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(graphics) == 0 {
		errMessage := "graphics not found"
		log.Warningf(ctx, "%v\n", errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	for i, graphic := range graphics {
		id := keys[i].StringID()
		filePath := "graphic/" + id + "." + graphic.FileType
		fileType := "" // this is unnecessary when GET request
		bucketName := infrastructure.MaterialBucketName(ctx)
		url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)

		if err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		graphics[i].ID = id
		graphics[i].URL = url

		replacedURL := strings.Replace(graphicThumbnailURL, "{id}", id, 1)
		graphics[i].ThumbnailURL = strings.Replace(replacedURL, "{fileType}", graphic.FileType, 1)
	}

	return c.JSON(http.StatusOK, graphics)
}
