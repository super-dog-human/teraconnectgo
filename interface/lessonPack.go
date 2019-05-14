package interface

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"archive/zip"
	"bytes"
	"cloudHelper"
	"io"
	"lessonType"
	"net/http"
	"utility"
)

// UpdateLessonPack is update lesson function.
func UpdateLessonPack(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	id := c.Param("id")

	ids := []string{id}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	var err error

	lesson := new(lessonType.Lesson)
	lesson.ID = id
	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	if err = datastore.Get(ctx, key, lesson); err != nil && err != datastore.ErrNoSuchEntity {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if lesson.IsPacked { // TODO remove me when end of beta.
		log.Warningf(ctx, "%+v\n", "already packed lesson")
		return c.JSON(http.StatusOK, "success.")
	}

	var graphicFileTypes map[string]string
	if graphicFileTypes, err = fetchGraphicFileTypesFromGCD(ctx, lesson.GraphicIDs); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = importGraphicsToZip(ctx, lesson.GraphicIDs, graphicFileTypes, zipWriter); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var lessonVoiceTexts []lessonType.LessonVoiceText
	query := datastore.NewQuery("LessonVoiceText").Filter("LessonID =", id)
	if _, err = query.GetAll(ctx, &lessonVoiceTexts); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = importVoiceToZip(ctx, lessonVoiceTexts, id, zipWriter); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = importLessonJsonToZip(ctx, id, zipWriter); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = removeUsedFilesInGCS(ctx, id, lessonVoiceTexts); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err = updateLessonAfterPacked(ctx, id); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	zipWriter.Close()

	zipFilePath := "lesson/" + id + ".zip"
	contentType := "application/zip"
	bucketName := utility.MaterialBucketName(ctx)
	if err := cloudHelper.CreateObjectToGCS(ctx, bucketName, zipFilePath, contentType, zipBuffer.Bytes()); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "success")
}

func importGraphicsToZip(ctx context.Context, usedGraphicIDs []string, graphicFileTypes map[string]string, zipWriter *zip.Writer) error {
	for _, graphicID := range usedGraphicIDs {
		fileType := graphicFileTypes[graphicID]
		filePathInGCS := "graphic/" + graphicID + "." + fileType
		bucketName := utility.MaterialBucketName(ctx)

		objectBytes, err := cloudHelper.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
		if err != nil {
			return err
		}

		filePathInZip := "graphics/" + graphicID + "." + fileType
		var f io.Writer
		f, err = zipWriter.Create(filePathInZip)
		if err != nil {
			return err
		}

		if _, err = f.Write(objectBytes); err != nil {
			return err
		}
	}

	return nil
}

func importVoiceToZip(ctx context.Context, voiceTexts []lessonType.LessonVoiceText, id string, zipWriter *zip.Writer) error {
	for _, voiceText := range voiceTexts {
		filePathInGCS := "voice/" + id + "/" + voiceText.FileID + ".ogg"
		bucketName := utility.MaterialBucketName(ctx)

		objectBytes, err := cloudHelper.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
		if err != nil {
			return err
		}

		filePathInZip := "voices/" + voiceText.FileID + ".ogg"
		var f io.Writer
		f, err = zipWriter.Create(filePathInZip)
		if err != nil {
			return err
		}

		if _, err = f.Write(objectBytes); err != nil {
			return err
		}
	}

	return nil
}

func importLessonJsonToZip(ctx context.Context, id string, zipWriter *zip.Writer) error {
	filePathInGCS := "lesson/" + id + ".json"
	bucketName := utility.MaterialBucketName(ctx)
	jsonBytes, err := cloudHelper.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
	if err != nil {
		return err
	}

	filePathInZip := "lesson.json"
	var f io.Writer
	f, err = zipWriter.Create(filePathInZip)
	if err != nil {
		return err
	}

	if _, err = f.Write(jsonBytes); err != nil {
		return err
	}

	return nil
}

func fetchGraphicFileTypesFromGCD(ctx context.Context, graphicIDs []string) (map[string]string, error) {
	var keys []*datastore.Key
	for _, id := range graphicIDs {
		keys = append(keys, datastore.NewKey(ctx, "Graphic", id, 0, nil))
	}

	graphicFileTypes := map[string]string{}
	graphics := make([]lessonType.Graphic, len(graphicIDs))
	if err := datastore.GetMulti(ctx, keys, graphics); err != nil {
		return nil, err
	} else {
		for i, g := range graphics {
			id := graphicIDs[i]
			graphicFileTypes[id] = g.FileType
		}
	}

	return graphicFileTypes, nil
}

func removeUsedFilesInGCS(ctx context.Context, id string, voiceTexts []lessonType.LessonVoiceText) error {
	var err error

	rawVoiceBucketName := utility.RawVoiceBucketName(ctx)
	voiceForTranscriptionBucketName := utility.VoiceForTranscriptionBucketName(ctx)
	for _, voiceText := range voiceTexts {
		filePathInGCS := id + "-" + voiceText.FileID + ".wav"

		err = cloudHelper.DeleteObjectsFromGCS(ctx, rawVoiceBucketName, filePathInGCS)
		if err != nil && err != storage.ErrObjectNotExist {
			return err
		}

		err = cloudHelper.DeleteObjectsFromGCS(ctx, voiceForTranscriptionBucketName, filePathInGCS)
		if err != nil && err != storage.ErrObjectNotExist {
			return err
		}
	}

	return nil
}

func updateLessonAfterPacked(ctx context.Context, id string) error {
	key := datastore.NewKey(ctx, "Lesson", id, 0, nil)
	lesson := new(lessonType.Lesson)
	lesson.ID = id

	var err error
	if err = datastore.Get(ctx, key, lesson); err != nil {
		return err
	}

	lesson.IsPacked = true
	if _, err = datastore.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}
