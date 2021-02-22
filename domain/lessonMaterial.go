package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonMaterial struct {
	ID                int64           `json:"id" datastore:"-"`
	Version           uint            `json:"version" datastore:"-"` // 同じLessonを親に持つもののCreatedの昇順をバージョンとする
	UserID            int64           `json:"userID"`
	AvatarID          int64           `json:"avatarID"`
	DurationSec       float32         `json:"durationSec" datastore:",noindex"`
	AvatarLightColor  string          `json:"avatarLightColor" datastore:",noindex"`
	BackgroundImageID int64           `json:"backgroundImageID"`
	BackgroundMusicID int64           `json:"backgroundMusicID"`
	Avatars           []LessonAvatar  `json:"avatars"`
	Graphics          []LessonGraphic `json:"graphics"`
	Drawings          []LessonDrawing `json:"drawings"`
	Speeches          []LessonSpeech  `json:"speeches"`
	Created           time.Time       `json:"created"`
	Updated           time.Time       `json:"updated"`
}

type LessonAvatar struct {
	Elapsedtime float32    `json:"elapsedtime"`
	DurationSec float32    `json:"durationSec"`
	Moving      Position3D `json:"moving,omitempty"`
}

type LessonGraphic struct {
	Elapsedtime float64 `json:"elapsedtime"`
	GraphicID   int64   `json:"graphicID"`
	Action      string  `json:"action"`
}

type LessonDrawing struct {
	Elapsedtime float32             `json:"elapsedtime"`
	DurationSec float32             `json:"durationSec"`
	Action      string              `json:"action"`
	Stroke      LessonDrawingStroke `json:"strokes"`
}

type LessonDrawingStroke struct {
	Clear     bool         `json:"clear"`
	Eraser    bool         `json:"eraser"`
	Width     int32        `json:"width,omitempty"`
	Height    int32        `json:"height,omitempty"`
	Color     string       `json:"color,omitempty"`
	LineWidth int32        `json:"lineWidth,omitempty"`
	Positions []Position2D `json:"positions,omitempty"`
}

type LessonSpeech struct {
	Elapsedtime float32  `json:"elapsedtime"`
	DurationSec float32  `json:"durationSec"`
	VoiceID     int64    `json:"voiceID"`
	Subtitle    Subtitle `json:"subtitle"`
	Caption     Caption  `json:"caption"`
}

type Subtitle struct {
	Body string `json:"body"`
}

type Caption struct {
	SizeVW          uint8  `json:"sizeVW"`
	Body            string `json:"body"`
	BodyColor       string `json:"bodyColor"`
	BorderColor     string `json:"borderColor"`
	HorizontalAlign string `json:"horizontalAlign"`
	VerticalAlign   string `json:"verticalAlign"`
}

func GetLessonMaterial(ctx context.Context, lessonID int64, lessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	query := datastore.NewQuery("LessonMaterial").Ancestor(ancestor).Order("-Created") // 降順
	var lessonMaterials []LessonMaterial
	keys, err := client.GetAll(ctx, query, &lessonMaterials)
	if err != nil {
		return err
	}

	if len(lessonMaterials) > 0 {
		*lessonMaterial = lessonMaterials[0]
		lessonMaterial.ID = keys[0].ID
		lessonMaterial.Version = uint(len(lessonMaterials))
	}

	return nil
}

func CreateLessonMaterial(ctx context.Context, lessonID int64, lessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	currentTime := time.Now()
	lessonMaterial.Created = currentTime
	lessonMaterial.Updated = currentTime

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IncompleteKey("LessonMaterial", ancestor)
	if _, err := client.Put(ctx, key, lessonMaterial); err != nil {
		return err
	}

	return nil
}

func UpdateLessonMaterial(ctx context.Context, id int64, lessonID int64, newLessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var lessonMaterial LessonMaterial

		ancestor := datastore.IDKey("Lesson", lessonID, nil)
		key := datastore.IDKey("LessonMaterial", id, ancestor)
		if err := tx.Get(key, lessonMaterial); err != nil {
			return nil
		}

		copier.Copy(&lessonMaterial, &newLessonMaterial)
		lessonMaterial.Updated = time.Now()

		if _, err := tx.Put(key, lessonMaterial); err != nil {
			return err
		}

		return nil
	})

	return nil
}
