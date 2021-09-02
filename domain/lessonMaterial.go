package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonMaterial struct {
	ID                   int64                `json:"id" datastore:"-"`
	UserID               int64                `json:"userID"`
	AvatarID             int64                `json:"avatarID"`
	Avatar               Avatar               `json:"avatar" datastore:"-"`
	AvatarLightColor     string               `json:"avatarLightColor" datastore:",noindex"`
	DurationSec          float32              `json:"durationSec" datastore:",noindex"`
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
	ElapsedTime float32   `json:"elapsedTime"`
	DurationSec float32   `json:"durationSec"`
	Positions   []float32 `json:"positions"` // x, y, zの順で格納
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
	StartAtSec  int32           `json:"startAtSec"`
	ServiceName string          `json:"serviceName"`
}

type LessonGraphic struct {
	ElapsedTime float32       `json:"elapsedTime"`
	GraphicID   int64         `json:"graphicID"`
	Action      GraphicAction `json:"action"`
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
	VoiceFileKey    string               `json:"voiceFileKey"`
	Subtitle        string               `json:"subtitle"`
	Caption         Caption              `json:"caption"`
	IsSynthesis     bool                 `json:"isSynthesis"`
	SynthesisConfig VoiceSynthesisConfig `json:"synthesisConfig"`
}

type Caption struct {
	Body            string `json:"body,omitempty"`
	BodyColor       string `json:"bodyColor,omitempty"`
	BorderColor     string `json:"borderColor,omitempty"`
	HorizontalAlign string `json:"horizontalAlign,omitempty"`
	VerticalAlign   string `json:"verticalAlign,omitempty"`
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

func UpdateLessonMaterial(ctx context.Context, id int64, lessonID int64, jsonBody *map[string]interface{}, targetFields *[]string) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		currentTime := time.Now()
		if _, err := updateLessonMaterialInTransaction(tx, id, lessonID, jsonBody, targetFields, currentTime); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func updateLessonMaterialInTransaction(tx *datastore.Transaction, id int64, lessonID int64, jsonBody *map[string]interface{}, targetFields *[]string, currentTime time.Time) (LessonMaterial, error) {
	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IDKey("LessonMaterial", id, ancestor)
	lessonMaterial := new(LessonMaterial)
	if err := tx.Get(key, lessonMaterial); err != nil {
		return *lessonMaterial, err
	}

	MergeJsonToStruct(jsonBody, lessonMaterial, targetFields)
	lessonMaterial.Updated = currentTime

	if _, err := tx.Put(key, lessonMaterial); err != nil {
		return *lessonMaterial, err
	}

	return *lessonMaterial, nil
}
