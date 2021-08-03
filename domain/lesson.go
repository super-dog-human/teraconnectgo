package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Lesson is the lesson infomation type.
type Lesson struct {
	ID                   int64             `json:"id" datastore:"-"`
	UserID               int64             `json:"userID"`
	MaterialID           int64             `json:"materialID"`
	PrevLessonID         int64             `json:"prevLessonID"`
	PrevLessonTitle      string            `json:"prevLessonTitle" datastore:"-"`
	NextLessonID         int64             `json:"nextLessonID"`
	NextLessonTitle      string            `json:"nextLessonTitle" datastore:"-"`
	NeedsRecording       bool              `json:"needsRecording"` // 収録画面での収録必要の有無
	IsEdited             bool              `json:"isEdited"`       // 編集画面から保存されたことがある
	IsIntroduction       bool              `json:"isIntroduction"` // 自己紹介用の授業
	IsPacked             bool              `json:"isPacked"`       // 公開準備の完了
	HasThumbnail         bool              `json:"hasThumbnail"`
	ThumbnailURL         string            `json:"thumbnailURL" datastore:"-"`
	AudioURL             string            `json:"audioURL" datastore:"-"`
	ZipURL               string            `json:"zipURL" datastore:"-"`
	Status               LessonStatus      `json:"status"`
	References           []LessonReference `json:"references"`
	Reviews              []LessonReview    `json:"reviews"`
	SubjectID            int64             `json:"subjectID"`
	SubjectName          string            `json:"subjectName"`
	JapaneseCategoryID   int64             `json:"japaneseCategoryID"`
	JapaneseCategoryName string            `json:"japaneseCategoryName"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	DurationSec          float64           `json:"durationSec"`
	ViewCount            int64             `json:"viewCount"`
	ViewKey              string            `json:"-"`
	SizeInBytes          int64             `json:"sizeInBytes"`
	Created              time.Time         `json:"created"`
	Updated              time.Time         `json:"updated"`
	Published            time.Time         `json:"published"`
}

// LessonReference is link to another web page.
type LessonReference struct {
	Name string `json:"name"`
	Isbn string `json:"isbn"` // ISBN13を想定
}

// LessonReview is review status of lesson by other users.
type LessonReview struct {
	ReviewerUserID int64     `json:"userID"`
	Comment        string    `json:"comment"`
	Created        time.Time `json:"created"`
}

const queueID string = "zipLesson"

func GetLessonByID(ctx context.Context, id int64) (Lesson, error) {
	lesson := new(Lesson)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *lesson, err
	}

	key := datastore.IDKey("Lesson", id, nil)
	if err := client.Get(ctx, key, lesson); err != nil {
		return *lesson, err
	}

	lesson.ID = id

	if err = SetLessonThumbnailURL(ctx, lesson); err != nil {
		return *lesson, err
	}

	return *lesson, nil
}

func GetLessonsByUserID(ctx context.Context, userID int64) ([]Lesson, error) {
	var lessons []Lesson

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("Lesson").Filter("UserID =", userID).Order("-Created")
	keys, err := client.GetAll(ctx, query, &lessons)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		lessons[i].ID = key.ID
		if err := SetLessonThumbnailURL(ctx, &lessons[i]); err != nil {
			return nil, err
		}
	}

	return lessons, nil
}

func CreateLesson(ctx context.Context, lesson *Lesson) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	if err := setCategoryAndSubject(ctx, lesson); err != nil {
		return err
	}

	currentTime := time.Now()
	lesson.Status = LessonStatusDraft
	lesson.Created = currentTime
	lesson.Updated = currentTime

	key, err := client.Put(ctx, datastore.IncompleteKey("Lesson", nil), lesson)

	if err != nil {
		return err
	}

	lesson.ID = key.ID

	return nil
}

func UpdateLesson(ctx context.Context, lesson *Lesson) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	lesson.Updated = time.Now()

	key := datastore.IDKey("Lesson", lesson.ID, nil)
	if _, err := client.Put(ctx, key, lesson); err != nil {
		return err
	}

	return nil
}

// UpdateLessonAndMaterialは、jsonのフィールドを既存のLesson/LessonMaterialへマージし、トランザクション中で二つのエンティティを更新します。
// jsonのフィールド名がlessonFieldsまたはlessonMaterialFieldsに含まれない場合、そのフィールドは無視されます。
func UpdateLessonAndMaterial(ctx context.Context, lesson *Lesson, needsCopyThumbnail bool, jsonBody *map[string]interface{}, lessonFields *[]string, lessonMaterialFields *[]string) error {
	currentStatus := lesson.Status
	currentSubjectID := lesson.SubjectID
	currentJapaneseCategoryID := lesson.JapaneseCategoryID

	MergeJsonToStruct(jsonBody, lesson, lessonFields)

	if lesson.SubjectID != currentSubjectID || lesson.JapaneseCategoryID != currentJapaneseCategoryID {
		if err := setCategoryAndSubject(ctx, lesson); err != nil {
			return err
		}
	}

	currentTime := time.Now()
	needsZipLesson := lesson.Status != LessonStatusDraft
	if needsZipLesson {
		lesson.Published = currentTime
	}

	if currentStatus != lesson.Status && needsCopyThumbnail {
		if err := CopyLessonThumbnail(ctx, lesson.ID, currentStatus, lesson.Status); err != nil {
			return err
		}
	}

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err = updateLessonInTransaction(tx, lesson); err != nil {
			return err
		}

		if err = updateLessonMaterialInTransaction(tx, lesson.MaterialID, lesson.ID, jsonBody, lessonMaterialFields); err != nil {
			return err
		}

		return nil
	})

	if needsZipLesson {
		taskName := strconv.FormatInt(lesson.ID, 10) + "-" + strconv.FormatInt(lesson.Published.UnixNano(), 10)
		taskEta := lesson.Published.Add(5 * time.Minute)
		taskBody := strconv.FormatInt(lesson.MaterialID, 10)
		if _, err := infrastructure.CreateTask(ctx, queueID, taskName, taskEta, taskBody); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func setCategoryAndSubject(ctx context.Context, lesson *Lesson) error {
	subject, err := GetSubject(ctx, lesson.SubjectID)
	if err != nil {
		return err
	}

	category, err := GetJapaneseCategory(ctx, lesson.JapaneseCategoryID, lesson.SubjectID)
	if err != nil {
		return err
	}

	lesson.SubjectName = subject.JapaneseName
	lesson.JapaneseCategoryName = category.Name

	return nil
}

func updateLessonInTransaction(tx *datastore.Transaction, lesson *Lesson) error {
	key := datastore.IDKey("Lesson", lesson.ID, nil)
	if _, err := tx.Put(key, lesson); err != nil {
		return err
	}

	return nil
}
