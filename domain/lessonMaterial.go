package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/imdario/mergo"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonMaterial struct {
	ID                   int64                `json:"id" datastore:"-"`
	UserID               int64                `json:"userID"`
	AvatarID             int64                `json:"avatarID"`
	Avatar               Avatar               `json:"avatar" datastore:"-"`
	DurationSec          float32              `json:"durationSec" datastore:",noindex"`
	AvatarLightColor     string               `json:"avatarLightColor" datastore:",noindex"`
	BackgroundImageID    int64                `json:"backgroundImageID"`
	BackgroundImageURL   string               `json:"backgroundImageURL" datastore:"-"`
	VoiceSynthesisConfig VoiceSynthesisConfig `json:"voiceSynthesisConfig" datastore:",noindex"`
	Avatars              []LessonAvatar       `json:"avatars" datastore:",noindex"`
	Graphics             []LessonGraphic      `json:"graphics" datastore:",noindex"`
	Drawings             []LessonDrawing      `json:"drawings" datastore:",noindex"`
	Embeddings           []LessonEmbedding    `json:"embeddings" datastore:",noindex"`
	Musics               []LessonMusic        `json:"musics" datastore:",noindex"`
	Speeches             []LessonSpeech       `json:"speeches" datastore:",noindex"`
	Created              time.Time            `json:"created" datastore:",noindex"`
	Updated              time.Time            `json:"updated" datastore:",noindex"`
}

type LessonAvatar struct {
	ElapsedTime float32    `json:"elapsedTime"`
	DurationSec float32    `json:"durationSec"`
	Moving      Position3D `json:"moving,omitempty"`
}

type LessonDrawing struct {
	ElapsedTime float32             `json:"elapsedTime"`
	DurationSec float32             `json:"durationSec"`
	Action      DrawingAction       `json:"action"`
	Units       []LessonDrawingUnit `json:"units"`
}

type LessonDrawingUnit struct {
	ElapsedTime float32             `json:"elapsedTime"`
	DurationSec float32             `json:"durationSec"`
	Action      DrawingUnitAction   `json:"action"`
	Stroke      LessonDrawingStroke `json:"stroke"`
}

type LessonDrawingStroke struct {
	Eraser    bool         `json:"eraser,omitempty"`
	Color     string       `json:"color,omitempty"`
	LineWidth int32        `json:"lineWidth,omitempty"`
	Positions []Position2D `json:"positions,omitempty"`
}

type LessonEmbedding struct {
	ElapsedTime float32         `json:"elapsedTime"`
	Action      EmbeddingAction `json:"action"`
	ContentID   string          `json:"contentID"`
	ServiceName string          `json:"type"`
}

type LessonGraphic struct {
	ElapsedTime float64        `json:"elapsedTime"`
	GraphicID   int64          `json:"graphicID"`
	Action      GraphicActrion `json:"action"`
}

type LessonMusic struct {
	ElapsedTime       float32     `json:"elapsedTime"`
	Action            MusicAction `json:"action"`
	BackgroundMusicID int64       `json:"backgroundMusicID"`
	Volume            float32     `json:"volume"`
	IsFading          bool        `json:"isFading"`
	IsLoop            bool        `json:"isLoop"`
}

type LessonSpeech struct {
	ElapsedTime     float32              `json:"elapsedTime"`
	DurationSec     float32              `json:"durationSec"`
	VoiceID         int64                `json:"voiceID"`
	Subtitle        string               `json:"subtitle"`
	Caption         Caption              `json:"caption"`
	IsSynthesis     bool                 `json:"isSynthesis"`
	SynthesisConfig VoiceSynthesisConfig `json:"synthesisConfig"`
}

type Caption struct {
	SizeVW          int8   `json:"sizeVW"`
	Body            string `json:"body"`
	BodyColor       string `json:"bodyColor"`
	BorderColor     string `json:"borderColor"`
	HorizontalAlign string `json:"horizontalAlign"`
	VerticalAlign   string `json:"verticalAlign"`
}

func GetLessonMaterial(ctx context.Context, id int64, lessonID int64, lessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IDKey("LessonMaterial", id, ancestor)
	if err := client.Get(ctx, key, lessonMaterial); err != nil {
		return err
	}

	lessonMaterial.ID = id
	lessonMaterial.BackgroundImageURL = infrastructure.GetPublicBackgroundImageURL(strconv.FormatInt(lessonMaterial.BackgroundImageID, 10))

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
	putKey, err := client.Put(ctx, key, lessonMaterial)
	if err != nil {
		return err
	}

	lessonMaterial.ID = putKey.ID

	return nil
}

func UpdateLessonMaterial(ctx context.Context, id int64, lessonID int64, newLessonMaterial *LessonMaterial) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		return UpdateLessonMaterialInTransaction(tx, id, lessonID, newLessonMaterial)
	})

	if err != nil {
		return err
	}

	return nil
}

func UpdateLessonMaterialInTransaction(tx *datastore.Transaction, id int64, lessonID int64, newLessonMaterial *LessonMaterial) error {
	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IDKey("LessonMaterial", id, ancestor)
	var lessonMaterial LessonMaterial
	if err := tx.Get(key, &lessonMaterial); err != nil {
		return err
	}

	if err := mergo.Merge(newLessonMaterial, lessonMaterial); err != nil {
		return err
	}

	newLessonMaterial.Updated = time.Now()

	if _, err := tx.Put(key, newLessonMaterial); err != nil {
		return err
	}

	return nil
}
