package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
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

func CreateLessonThumbnailBlankFile(ctx context.Context, id int64, isPublic bool) (string, error) {
	fileName := "thumbnail"
	fileRequest := infrastructure.FileRequest{
		Extension:   "png",
		ContentType: "image/png",
	}

	fileDir := "lesson/" + strconv.FormatInt(id, 10)
	var url string
	var err error

	if isPublic {
		url, err = infrastructure.CreateBlankFileToPublicGCS(ctx, fileName, fileDir, fileRequest)
		if err != nil {
			return "", err
		}
	} else {
		url, err = infrastructure.CreateBlankFileToGCS(ctx, fileName, fileDir, fileRequest)
		if err != nil {
			return "", err
		}
	}

	return url, nil
}

func CopyLessonThumbnail(ctx context.Context, id int64, currentStatus LessonStatus, newStatus LessonStatus) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	var srcBucket string
	var destBucket string
	if currentStatus == LessonStatusPublic {
		srcBucket = infrastructure.PublicBucketName()
	} else {
		srcBucket = infrastructure.MaterialBucketName()
	}

	if newStatus == LessonStatusPublic {
		destBucket = infrastructure.PublicBucketName()
	} else {
		destBucket = infrastructure.MaterialBucketName()
	}

	if srcBucket == destBucket {
		return nil
	}

	fileID := strconv.FormatInt(id, 10)
	objectPath := thumbnailFilePath(fileID)
	src := client.Bucket(srcBucket).Object(objectPath)
	dst := client.Bucket(destBucket).Object(objectPath)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}

	return nil
}

func thumbnailFilePath(id string) string {
	return "lesson/" + id + "/thumbnail.png"
}

func createPublicURL(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	return infrastructure.CloudStorageURL + infrastructure.PublicBucketName() + "/" + thumbnailFilePath(idStr)
}

func createSignedURL(ctx context.Context, id int64) (string, error) {
	idStr := strconv.FormatInt(id, 10)
	fileType := "" // this is unnecessary when GET request
	bucketName := infrastructure.MaterialBucketName()
	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, thumbnailFilePath(idStr), "GET", fileType)

	if err != nil {
		return "", err
	}

	return url, nil
}
