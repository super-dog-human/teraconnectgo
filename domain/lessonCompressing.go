package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func createLessonMaterialForCompressing(ctx context.Context, id string, lessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("LessonMaterialForCompressing", id, nil)
	if _, err := client.Put(ctx, key, lessonMaterial); err != nil {
		return err
	}

	return nil
}

func createCompressingTask(ctx context.Context, taskName string, materialID int64, currentTime time.Time) error {
	var taskEta time.Time
	if infrastructure.AppEnv() == "production" {
		taskEta = currentTime.Add(5 * time.Minute)
	} else {
		taskEta = currentTime.Add(1 * time.Minute)
	}
	// タスクに必要な情報はtaskNameで事足りるのでmessageは空文字でよい
	if _, err := infrastructure.CreateTask(ctx, taskName, taskEta, ""); err != nil {
		return err
	}

	return nil
}
