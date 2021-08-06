package domain

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonCompressingErrorCode uint

const (
	AlreadyCompressed  LessonCompressingErrorCode = 1
	AnotherTaskWillRun LessonCompressingErrorCode = 2
)

func (e LessonCompressingErrorCode) Error() string {
	switch e {
	case AlreadyCompressed:
		return "lesson compressing is already executed."
	case AnotherTaskWillRun:
		return "created another task after this."
	default:
		return "unknown error"
	}
}

type LessonMaterialForCompressing struct {
	LessonMaterial
	IsCompressing bool
}

func GetLessonMaterialForCompressing(ctx context.Context, id string) (LessonMaterialForCompressing, error) {
	key := datastore.NameKey("LessonMaterialForCompressing", id, nil)
	lessonMaterial := new(LessonMaterialForCompressing)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *lessonMaterial, err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(key, lessonMaterial); err != nil {
			return err
		}

		if lessonMaterial.IsCompressing {
			return AlreadyCompressed
		}

		lessonMaterial.IsCompressing = true
		if _, err := tx.Put(key, lessonMaterial); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return *lessonMaterial, err
	}

	return *lessonMaterial, nil
}

func UpdateLessonAfterCompressing(ctx context.Context, id int64, updated time.Time) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.IDKey("Lesson", id, nil)
		lesson := new(Lesson)
		if err := tx.Get(key, lesson); err != nil {
			return err
		}

		if lesson.Updated.After(updated) {
			return AnotherTaskWillRun
		}

		lesson.Published = updated
		if _, err := tx.Put(key, lesson); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteLessonMaterialForCompress(ctx context.Context, id string) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("LessonMaterialForCompressing", id, nil)
	if err := client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}

func CompressLesson(ctx context.Context, lessonID int64, taskName string, lessonMaterial *LessonMaterialForCompressing) error {
	workingDir := "/tmp/" + taskName

	if err := os.Mkdir(workingDir, 0644); err != nil {
		return err
	}

	if err := downloadOrCreateVoiceFiles(ctx, lessonID, workingDir, lessonMaterial); err != nil {
		return err
	}

	if err := mixAllAudios(lessonID, workingDir, &lessonMaterial.Speeches); err != nil {
		return err
	}

	if err := createCompressedMaterialToCloudStorage(lessonMaterial); err != nil {
		return err
	}

	if err := os.RemoveAll(workingDir); err != nil {
		return nil // 一時ファイルの削除に失敗しても実害はないので握り潰す
	}

	return nil
}

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
	taskEta := currentTime.Add(5 * time.Minute)
	// タスクに必要な情報はtaskNameで事足りるのでmessageは空文字でよい
	if _, err := infrastructure.CreateTask(ctx, taskName, taskEta, ""); err != nil {
		return err
	}

	return nil
}

func downloadOrCreateVoiceFiles(ctx context.Context, lessonID int64, workingDir string, lessonMaterial *LessonMaterialForCompressing) error {
	for _, speech := range lessonMaterial.Speeches {
		if !speech.IsSynthesis {
			continue
		}

		bucketName := infrastructure.MaterialBucketName()
		var voiceFile []byte
		var err error

		if speech.VoiceID == 0 && speech.Subtitle != "" {
			voiceParams := CreateSynthesisVoiceParam{LessonID: lessonID, Text: speech.Subtitle}
			if speech.SynthesisConfig.Name == "" {
				voiceParams.VoiceSynthesisConfig = lessonMaterial.VoiceSynthesisConfig
			} else {
				voiceParams.VoiceSynthesisConfig = speech.SynthesisConfig
			}
			// フロント側から声生成は済んでいるはずなので、ここで都度生成することはほぼないはず
			voiceFile, err = CreateSynthesizedVoice(ctx, &voiceParams, bucketName, "")
			if err != nil {
				return err
			}
		} else if speech.VoiceID != 0 {
			filePath := CloudStorageVoiceFilePath(lessonID, speech.VoiceID)
			voiceFile, err = infrastructure.GetFileFromGCS(ctx, bucketName, filePath)
			if err != nil {
				return err
			}
		}

		tmpFilePath := fmt.Sprintf("%d.mp3", speech.VoiceID)
		if err := ioutil.WriteFile(tmpFilePath, voiceFile, 0644); err != nil {
			return err
		}
	}

	return nil
}

func mixAllAudios(lessonID int64, workindDir string, speeches *[]LessonSpeech) error {
	return nil
}

func createCompressedMaterialToCloudStorage(lessonMaterial *LessonMaterialForCompressing) error {
	// jsonをzstd圧縮する
	return nil
}
