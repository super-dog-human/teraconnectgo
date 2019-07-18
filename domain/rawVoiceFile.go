package domain

import (
	"context"

	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
)

const contentType = "audio/wav"

func CreateBlankRawVoiceFile(ctx context.Context, lessonID string, fileID string) error {
	bucketName := infrastructure.RawVoiceBucketName(ctx)
	filePath := filePath(lessonID, fileID)
	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, filePath, contentType, nil); err != nil {
		return err
	}

	return nil
}

func GetRawVoiceFileSignedURLForUpload(ctx context.Context, lessonID string, fileID string) (string, error) {
	bucketName := infrastructure.RawVoiceBucketName(ctx)
	filePath := filePath(lessonID, fileID)
	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "PUT", contentType)
	if err != nil {
		return "", err
	}

	return url, nil
}

func filePath(lessonID string, fileID string) string {
	return lessonID + "-" + fileID + ".wav"
}