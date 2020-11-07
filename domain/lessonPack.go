package domain

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func CreateLessonZip(ctx context.Context, lesson Lesson, graphicFileTypes map[int64]string, voiceTexts []RawVoiceText) (*bytes.Buffer, error) {
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

func UploadLessonZipToGCS(ctx context.Context, lessonID int64, zip *bytes.Buffer) error {
	zipFilePath := fmt.Sprintf("lesson/%d.zip", lessonID)
	contentType := "application/zip"
	bucketName := infrastructure.MaterialBucketName()
	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, zipFilePath, contentType, zip.Bytes()); err != nil {
		return err
	}

	return nil
}

func RemoveUsedFilesInGCS(ctx context.Context, id int64, voiceTexts []RawVoiceText) error {
	var err error

	rawVoiceBucketName := infrastructure.RawVoiceBucketName()
	voiceForTranscriptionBucketName := infrastructure.VoiceForTranscriptionBucketName()
	for _, voiceText := range voiceTexts {
		filePathInGCS := fmt.Sprintf("%d-%s.wav", id, voiceText.FileID)

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

func addGraphicsToZip(ctx context.Context, usedGraphicIDs []int64, graphicFileTypes map[int64]string, zipWriter *zip.Writer) error {
	for _, graphicID := range usedGraphicIDs {
		fileType := graphicFileTypes[graphicID]
		filePathInGCS := fmt.Sprintf("graphic/%d.%s", graphicID, fileType)
		bucketName := infrastructure.MaterialBucketName()

		objectBytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
		if err != nil {
			return err
		}

		filePathInZip := fmt.Sprintf("graphics/%d.%s", graphicID, fileType)
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

func addVoiceToZip(ctx context.Context, voiceTexts []RawVoiceText, id int64, zipWriter *zip.Writer) error {
	for _, voiceText := range voiceTexts {
		filePathInGCS := fmt.Sprintf("voice/%d/%s.ogg", id, voiceText.FileID)
		bucketName := infrastructure.MaterialBucketName()

		objectBytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, filePathInGCS)
		if err != nil {
			return err
		}

		filePathInZip := fmt.Sprintf("voices/%d/%s.ogg", id, voiceText.FileID)
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

func addLessonJSONToZip(ctx context.Context, id int64, zipWriter *zip.Writer) error {
	filePathInGCS := fmt.Sprintf("lesson/%d.json", id)
	bucketName := infrastructure.MaterialBucketName()
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
