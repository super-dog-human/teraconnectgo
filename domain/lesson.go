package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Lesson is the lesson infomation type.
type Lesson struct {
	ID             int64     `json:"id" datastore:"-"`
	SubjectName    string    `json:"subjectName"`
	CategoryName   string    `json:"categoryName"`
	AvatarID       int64     `json:"avatarID"`
	Avatar         Avatar    `json:"avatar" datastore:"-"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	DurationSec    float64   `json:"durationSec"`
	ThumbnailURL   string    `json:"thumbnailURL" datastore:"-"`
	GraphicIDs     []int64   `json:"graphicIDs"`
	Graphics       []Graphic `json:"graphics" datastore:"-"`
	ViewCount      int64     `json:"viewCount"`
	Version        int64     `json:"version"`
	ViewKey        string    `json:"-"`
	IsIntroduction bool      `json:"isIntroduction"`
	IsPacked       bool      `json:"isPacked"`
	IsPublic       bool      `json:"isPublic"`
	UserID         int64     `json:"userID"`
	SizeInBytes    int64     `json:"sizeInBytes"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
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

	return *lesson, nil
}

func GetLessonsByUserID(ctx context.Context, userID int64) ([]Lesson, error) {
	var lessons []Lesson

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("Lesson").Filter("UserID =", userID)
	if _, err := client.GetAll(ctx, query, &lessons); err != nil {
		return nil, err
	}

	return lessons, nil
}

func CreateLesson(ctx context.Context, lesson *Lesson) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	currentTime := time.Now()
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

func DeleteLessonInTransactionByID(tx *datastore.Transaction, id int64) error {
	key := datastore.IDKey("Lesson", id, nil)
	if err := tx.Delete(key); err != nil {
		return err
	}

	return nil
}
