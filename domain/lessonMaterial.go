package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetLessonMaterialFromGCS(ctx context.Context, lessonID int64) (LessonMaterial, error) {
	lessonMaterial := new(LessonMaterial)

	bucketName := infrastructure.MaterialBucketName()
	bytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, lessonFilePath(lessonID))

	if err != nil {
		if err == storage.ErrObjectNotExist {
			return *lessonMaterial, err
		}
		return *lessonMaterial, err
	}

	if err := json.Unmarshal(bytes, lessonMaterial); err != nil {
		return *lessonMaterial, err
	}

	return *lessonMaterial, nil
}

func CreateLessonMaterialFileToGCS(ctx context.Context, lessonID int64, lessonMaterial LessonMaterial) error {
	contents, err := json.Marshal(lessonMaterial)
	if err != nil {
		return err
	}

	contentType := "application/json"
	bucketName := infrastructure.MaterialBucketName()
	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, lessonFilePath(lessonID), contentType, contents); err != nil {
		return err
	}

	return nil
}

func lessonFilePath(lessonID int64) string {
	return fmt.Sprintf("lesson/%d.json", lessonID)
}
