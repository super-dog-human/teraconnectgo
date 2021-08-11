package usecase

import (
	"errors"
	"log"
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// CompreesLesson
func CompreesLesson(request *http.Request, lessonID int64, taskName string, queuedUnixNanoTime int64) error {
	ctx := request.Context()

	lesson, err := domain.GetLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}

	if lesson.Status == domain.LessonStatusDraft {
		log.Println("LessonのStatusがDraftになっていたので終了")
		return nil
	}

	if lesson.Updated.UnixNano() > queuedUnixNanoTime {
		log.Println("LessonのUpdatedがタスク作成時刻より後のため終了")
		return nil
	}

	lessonMaterial, err := domain.GetLessonMaterialForCompressing(ctx, taskName)
	if err != nil {
		if ok := errors.Is(err, domain.AlreadyCompressed); ok {
			log.Println("LessonMaterialForCompressingのIsCompressingがtrueなので終了")
			return nil
		}
		log.Printf("GetLessonMaterialForCompressing error: %v\n", err.Error())
		return err
	}

	if err := domain.CompressLesson(ctx, &lesson, taskName, &lessonMaterial); err != nil {
		log.Printf("CompressLesson error: %v\n", err.Error())
		return err
	}

	if err = domain.UpdateLessonAfterCompressing(ctx, lessonID, lessonMaterial.DurationSec, lesson.Updated); err != nil {
		// 失敗時も実害はないのでそのまま終了
		if ok := errors.Is(err, domain.AnotherTaskWillRun); ok {
			log.Println("このタスク実行中に更新があったのでPublishedは更新せずに終了")
			return nil
		}
		log.Println(err.Error())
		log.Println("LessonのPublicの更新に失敗")
		return nil
	}

	if err = domain.DeleteLessonMaterialForCompress(ctx, taskName); err != nil {
		// 同上
		log.Println(err.Error())
		log.Println("LessonMaterialForCompressの削除に失敗")
		return nil
	}

	return nil
}
