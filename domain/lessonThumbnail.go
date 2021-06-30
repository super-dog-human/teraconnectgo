package domain

import (
	"context"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func SetLessonThumbnailURL(ctx context.Context, lesson *Lesson) error {
	if !lesson.HasThumbnail {
		return nil
	}

	if lesson.Status == LessonStatusPublic {
		lesson.ThumbnailURL = createPublicURL(lesson.ID)
	} else {
		url, err := createSignedURL(ctx, lesson.ID)
		if err != nil {
			return err
		}
		lesson.ThumbnailURL = url
	}

	return nil
}

func createPublicURL(id int64) string {
	fileID := strconv.FormatInt(id, 10)
	return "https://storage.googleapis.com/" + infrastructure.PublicBucketName() + "/lesson_thumbnail/" + fileID + ".png"
}

func createSignedURL(ctx context.Context, id int64) (string, error) {
	fileID := strconv.FormatInt(id, 10)
	filePath := "lesson_thumbnail/" + fileID + ".png"
	fileType := "" // this is unnecessary when GET request
	bucketName := infrastructure.MaterialBucketName()
	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", fileType)

	if err != nil {
		return "", err
	}

	return url, nil
}
