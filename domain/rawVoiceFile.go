package domain

import (
	"context"
	"fmt"

	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

const contentType = "audio/wav"

func CreateBlankRawVoiceFile(ctx context.Context, lessonID int64, fileID string) error {
	bucketName := infrastructure.RawVoiceBucketName()
	filePath := voiceFilePath(lessonID, fileID)
	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, filePath, contentType, nil); err != nil {
		return err
	}

	return nil
}

func GetRawVoiceFileSignedURLForUpload(ctx context.Context, lessonID int64, fileID string) (string, error) {
	bucketName := infrastructure.RawVoiceBucketName()
	filePath := voiceFilePath(lessonID, fileID)
	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "PUT", contentType)
	if err != nil {
		return "", err
	}

	return url, nil
}

func voiceFilePath(lessonID int64, fileID string) string {
	return fmt.Sprintf("%d-%s.wav", lessonID, fileID)
}
