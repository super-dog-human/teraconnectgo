package interface

import (
	"cloudHelper"
	"net/http"
	"utility"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const contentType = "audio/wav"

// PostRawVoice is create blank wav file to Cloud Storage for direct upload from client.
func PostRawVoice(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())

	request := new(postRequest)
	if err := c.Bind(request); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lessonID := request.LessonID
	ids := []string{lessonID}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	fileID := xid.New().String()
	filePath := lessonID + "-" + fileID + ".wav"
	bucketName := utility.RawVoiceBucketName(ctx)

	if err := cloudHelper.CreateObjectToGCS(ctx, bucketName, filePath, contentType, nil); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	url, err := cloudHelper.GetGCSSignedURL(ctx, bucketName, filePath, "PUT", contentType)
	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURL{FileID: fileID, SignedURL: url})
}

type postRequest struct {
	LessonID string `json:"lesson_id"`
}

type signedURL struct {
	FileID    string `json:"file_id"`
	SignedURL string `json:"signed_url"`
}
