package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// GetStorageObjects is get signed URLs of files.
func GetStorageObjects(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	jsonString := c.Request().Header.Get("X-Get-Params")
	var fileRequests []fileRequest
	if err := json.Unmarshal([]byte(jsonString), &fileRequests); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	urlLength := len(fileRequests)
	urls := make([]signedURL, urlLength)

	for i, request := range fileRequests {
		// TODO check user permission
		// TODO check file exists

		filePath := filePath(request.Entity, request.ID, request.Extension)
		bucketName := infrastructure.MaterialBucketName(ctx)
		url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
		if err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		urls[i] = signedURL{request.ID, url}
	}

	return c.JSON(http.StatusOK, urlResponses{SignedURLs: urls})
}

// PostStorageObjects is create blank object to Cloud Storage for direct upload from client.
func PostStorageObjects(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	request := new(postStorageObjectRequest)
	if err := c.Bind(request); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	bucketName := infrastructure.MaterialBucketName(ctx)
	urls := make([]signedURL, len(request.FileRequests))

	for i, fileRequest := range request.FileRequests {
		fileID := xid.New().String()
		filePath := filePath(fileRequest.Entity, fileID, fileRequest.Extension)

		if err := infrastructure.CreateObjectToGCS(ctx, bucketName, filePath, fileRequest.ContentType, nil); err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "PUT", fileRequest.ContentType)
		if err != nil {
			log.Errorf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if fileRequest.Entity == "graphic" {
			graphic := new(domain.Graphic)
			graphic.Created = time.Now()
			graphic.FileType = fileRequest.Extension
			// graphic.UserID  = "foo"	// TODO
			key := datastore.NewKey(ctx, "Graphic", fileID, 0, nil)
			if _, err = datastore.Put(ctx, key, graphic); err != nil {
				log.Errorf(ctx, "%+v\n", errors.WithStack(err))
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		} else if fileRequest.Entity == "avatar" {
			avatar := new(domain.Avatar)
			avatar.Created = time.Now()
			// avatar.UserID  = "foo"	// TODO
			key := datastore.NewKey(ctx, "Avatar", fileID, 0, nil)
			if _, err = datastore.Put(ctx, key, avatar); err != nil {
				log.Errorf(ctx, "%+v\n", errors.WithStack(err))
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		urls[i] = signedURL{fileID, url}
	}

	return c.JSON(http.StatusOK, urlResponses{SignedURLs: urls})
}

func filePath(entity string, id string, extension string) string {
	return strings.ToLower(entity) + "/" + id + "." + extension
}

type postStorageObjectRequest struct {
	FileRequests []fileRequest `json:"fileRequests"`
}

type fileRequest struct {
	ID          string `json:"id"`
	Entity      string `json:"entity"`
	Extension   string `json:"extension"`
	ContentType string `json:"contentType"`
}

type urlResponses struct {
	SignedURLs []signedURL `json:"signedURLs"`
}

type signedURL struct {
	FileID    string `json:"fileID"`
	SignedURL string `json:"signedURL"`
}
