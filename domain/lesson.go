package domain

import (
	"context"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
	"google.golang.org/api/iterator"
)

// Lesson is the lesson infomation type.
type Lesson struct {
	ID                   int64             `json:"id" datastore:"-"`
	UserID               int64             `json:"userID"`
	Author               User              `json:"author" datastore:"-"`
	MaterialID           int64             `json:"materialID"`
	AvatarID             int64             `json:"avatarID"`                              // 公開処理完了時にLessonMaterialの値で更新される
	AvatarLightColor     string            `json:"avatarLightColor" datastore:",noindex"` // 公開処理完了時にLessonMaterialの値で更新される
	Avatar               Avatar            `json:"avatar,omitempty" datastore:"-"`
	PrevLessonID         int64             `json:"prevLessonID" datastore:",noindex"`
	PrevLessonTitle      string            `json:"prevLessonTitle" datastore:"-"`
	NextLessonID         int64             `json:"nextLessonID" datastore:",noindex"`
	NextLessonTitle      string            `json:"nextLessonTitle" datastore:"-"`
	NeedsRecording       bool              `json:"needsRecording" datastore:",noindex"` // 収録画面での収録必要の有無
	IsIntroduction       bool              `json:"isIntroduction"`                      // 自己紹介用の授業
	HasThumbnail         bool              `json:"hasThumbnail" datastore:",noindex"`
	ThumbnailURL         string            `json:"thumbnailURL" datastore:"-"`
	SpeechURL            string            `json:"speechURL" datastore:"-"`
	BodyURL              string            `json:"bodyURL" datastore:"-"`
	Status               LessonStatus      `json:"status"`
	References           []LessonReference `json:"references" datastore:",noindex"`
	Reviews              []LessonReview    `json:"reviews" datastore:",noindex"`
	SubjectID            int64             `json:"subjectID"`
	SubjectName          string            `json:"subjectName" datastore:",noindex"`
	JapaneseCategoryID   int64             `json:"japaneseCategoryID"`
	JapaneseCategoryName string            `json:"japaneseCategoryName" datastore:",noindex"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	DurationSec          float32           `json:"durationSec" datastore:",noindex"`
	ViewCount            int64             `json:"viewCount" datastore:",noindex"`
	ViewKey              string            `json:"viewKey" datastore:",noindex"`
	Version              int32             `json:"version" datastore:",noindex"`
	Created              time.Time         `json:"created" datastore:",noindex"`
	Updated              time.Time         `json:"updated" datastore:",noindex"`
	Published            time.Time         `json:"published"` // 公開処理完了時にLessonMaterialのUpdatedの値で更新される
}

type ShortLesson struct {
	ID          int64  `json:"id" datastore:"-"`
	UserID      int64  `json:"userID"`
	Title       string `json:"title"`
	Description string `json:"description"`
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

func GetLessonsByCategoryID(ctx context.Context, cursorStr string, categoryID int64) ([]ShortLesson, string, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, "", err
	}

	const lessonPageSize = 18
	query := datastore.NewQuery("Lesson").Project("UserID", "Title", "Description").
		Filter("JapaneseCategoryID =", categoryID).Filter("Status = ", int32(LessonStatusPublic)).Order("-Published").Limit(lessonPageSize)

	if cursorStr != "" {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return nil, "", err
		}
		query = query.Start(cursor)
	}

	var lessons []ShortLesson
	it := client.Run(ctx, query)
	for {
		var lesson ShortLesson
		key, err := it.Next(&lesson)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, "", err
		}
		lesson.ID = key.ID
		lessons = append(lessons, lesson)
	}

	nextCursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}

	return lessons, nextCursor.String(), nil
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

func CreateIntroductionLesson(ctx context.Context, user *User) (int64, error) {
	var lessonID int64

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return lessonID, err
	}

	query := datastore.NewQuery("Lesson").KeysOnly().Filter("UserID =", user.ID).Filter("IsIntroduction =", true).Limit(1)
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return lessonID, err
	}

	if len(keys) > 0 {
		return keys[0].ID, nil // 既に自己紹介授業が作られていればそれを使用する
	}

	lesson := Lesson{UserID: user.ID, Title: "はじめまして、" + user.Name + "です。", IsIntroduction: true, Created: time.Now()}
	key := datastore.IncompleteKey("Lesson", nil)
	putKey, err := client.Put(ctx, key, &lesson)
	if err != nil {
		return lessonID, err
	}

	return putKey.ID, nil
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
func UpdateLessonAndMaterial(ctx context.Context, lesson *Lesson, needsCopyThumbnail bool, requestID string, jsonBody *map[string]interface{}, lessonFields *[]string, lessonMaterialFields *[]string) error {
	currentStatus := lesson.Status
	currentSubjectID := lesson.SubjectID
	currentJapaneseCategoryID := lesson.JapaneseCategoryID

	MergeJsonToStruct(jsonBody, lesson, lessonFields)

	if lesson.SubjectID != currentSubjectID || lesson.JapaneseCategoryID != currentJapaneseCategoryID {
		if err := setCategoryAndSubject(ctx, lesson); err != nil {
			return err
		}
	}

	if currentStatus != lesson.Status && needsCopyThumbnail {
		if err := CopyLessonThumbnail(ctx, lesson.ID, currentStatus, lesson.Status); err != nil {
			return err
		}
	}

	currentTime := time.Now()
	if lesson.Status == LessonStatusLimited && lesson.ViewKey == "" {
		uuid, err := UUIDWithoutHypen()
		if err != nil {
			return err
		}
		lesson.ViewKey = uuid
	} else if lesson.Status == LessonStatusPublic && lesson.ViewKey != "" {
		lesson.ViewKey = ""
	}
	lesson.Updated = currentTime

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	var lessonMaterial LessonMaterial
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err = updateLessonInTransaction(tx, lesson); err != nil {
			return err
		}

		if lessonMaterial, err = updateLessonMaterialInTransaction(tx, lesson.MaterialID, lesson.ID, jsonBody, lessonMaterialFields, currentTime); err != nil {
			return err
		}

		return nil
	})

	if lesson.Status != LessonStatusDraft {
		taskName := infrastructure.LessonCompressingTaskName(lesson.ID, currentTime, requestID)
		if err := createLessonMaterialForCompressing(ctx, taskName, &lessonMaterial); err != nil {
			return err
		}

		if err := createCompressingTask(ctx, taskName, lesson.MaterialID, currentTime); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if currentStatus == LessonStatusPublic && lesson.Status != LessonStatusPublic {
		// 公開を取りやめる際は検索インデックスから登録を削除
		client := search.NewClient(os.Getenv("ALGOLIA_APPLICATION_ID"), os.Getenv("ALGOLIA_ADMIN_API_KEY"))
		index := client.InitIndex(os.Getenv("ALGOLIA_INDEX_NAME"))
		if _, err := index.DeleteObject(strconv.FormatInt(lesson.ID, 10)); err != nil {
			return err
		}
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
