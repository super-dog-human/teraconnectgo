package domain

import (
	"archive/zip"
	"bytes"
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
)

func CreateLessonZip(ctx context.Context, lesson Lesson, graphicFileTypes map[string]string, voiceTexts []RawVoiceText) (*bytes.Buffer, error) {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	var err error
	if err = addGraphicsToZip(ctx, lesson.GraphicIDs, graphicFileTypes, zipWriter); err != nil {
		return zipBuffer, err
	}

	if err = addVoiceToZip(ctx, voiceTexts, lesson.ID, zipWriter); err != nil {
		return zipBuffer, err
	}

	if err = addLessonJSONToZip(ctx, lesson.ID, zipWriter); err != nil {
		return zipBuffer, err
	}

	zipWriter.Close()

	return zipBuffer, nil
}

func UploadLessonZipToGCS(ctx context.Context, lessonID string, zip *bytes.Buffer) error {
	zipFilePath := "lesson/" + lessonID + ".zip"
	contentType := "application/zip"
	bucketName := infrastructure.MaterialBucketName(ctx)
	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, zipFilePath, contentType, zip.Bytes()); err != nil {
		return err
	}

	return nil
}

func RemoveUsedFilesInGCS(ctx context.Context, id string, voiceTexts []RawVoiceText) error {
	var err error

	rawVoiceBucketName := infrastructure.RawVoiceBucketName(ctx)
	voiceForTranscriptionBucketName := infrastructure.VoiceForTranscriptionBucketName(ctx)
	for _, voiceText := range voiceTexts {
		filePathInGCS := id + "-" + voiceText.FileID + ".wav"

		err = infrastructure.DeleteObjectsFromGCS(ctx, rawVoiceBucketName, filePathInGCS)
		if err != nil && err != storage.ErrObjectNotExist {
			return err
		}

		err = infrastructure.DeleteObjectsFromGCS(ctx, voiceForTranscriptionBucketName, filePathInGCS)
		if err != nil && err != storage.ErrObjectNotExist {
			return err
		}
	}

	return nil
}

func addGraphicsToZip(ctx context.Context, usedGraphicIDs []string, graphicFileTypes map[string]string, zipWriter *zip.Writer) error {
	for _, graphicID := range usedGraphicIDs {
		fileType := graphicFileTypes[graphicID]
		filePathInGCS := "graphic/" + graphicID + "." + fileType
		bucketName := infrastructure.MaterialBucketName(ctx)

		objectBytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
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

func addVoiceToZip(ctx context.Context, voiceTexts []RawVoiceText, id string, zipWriter *zip.Writer) error {
	for _, voiceText := range voiceTexts {
		filePathInGCS := "voice/" + id + "/" + voiceText.FileID + ".ogg"
		bucketName := infrastructure.MaterialBucketName(ctx)

		objectBytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
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

func addLessonJSONToZip(ctx context.Context, id string, zipWriter *zip.Writer) error {
	filePathInGCS := "lesson/" + id + ".json"
	bucketName := infrastructure.MaterialBucketName(ctx)
	jsonBytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
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
