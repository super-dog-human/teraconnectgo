package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonMaterial struct {
	ID                   int64                `json:"id" datastore:"-"`
	Version              uint                 `json:"version" datastore:"-"` // 同じLessonを親に持つもののCreatedの昇順をバージョンとする
	UserID               int64                `json:"userID"`
	AvatarID             int64                `json:"avatarID"`
	DurationSec          float32              `json:"durationSec" datastore:",noindex"`
	AvatarLightColor     string               `json:"avatarLightColor" datastore:",noindex"`
	BackgroundImageID    int64                `json:"backgroundImageID"`
	BackgroundImageURL   string               `json:"backgroundImageURL" datastore:"-"`
	VoiceSynthesisConfig VoiceSynthesisConfig `json:"voiceSynthesisConfig"`
	Avatars              []LessonAvatar       `json:"avatars"`
	Graphics             []LessonGraphic      `json:"graphics"`
	Drawings             []LessonDrawing      `json:"drawings"`
	Embeddings           []LessonEmbedding    `json:"embeddings"`
	Musics               []LessonMusic        `json:"musics"`
	Speeches             []LessonSpeech       `json:"speeches"`
	Created              time.Time            `json:"created"`
	Updated              time.Time            `json:"updated"`
}

type LessonAvatar struct {
	ElapsedTime float32    `json:"elapsedTime"`
	DurationSec float32    `json:"durationSec"`
	Moving      Position3D `json:"moving,omitempty"`
}

type LessonGraphic struct {
	ElapsedTime float64 `json:"elapsedTime"`
	GraphicID   int64   `json:"graphicID"`
	Action      string  `json:"action"`
}

type LessonDrawing struct {
	ElapsedTime float32             `json:"elapsedTime"`
	DurationSec float32             `json:"durationSec"`
	Action      string              `json:"action"` // draw/clear/show/hide
	Units       []LessonDrawingUnit `json:"units"`
}

type LessonEmbedding struct {
	ElapsedTime float32 `json:"elapsedTime"`
	Action      string  `json:"action"` // show/hide
	ContentID   string  `json:"contentID"`
	ServiceName string  `json:"type"`
}

type LessonDrawingUnit struct {
	ElapsedTime float32             `json:"elapsedTime"`
	DurationSec float32             `json:"durationSec"`
	Action      string              `json:"action"` //draw/undo
	Stroke      LessonDrawingStroke `json:"stroke"`
}

type LessonDrawingStroke struct {
	Eraser    bool         `json:"eraser,omitempty"`
	Color     string       `json:"color,omitempty"`
	LineWidth int32        `json:"lineWidth,omitempty"`
	Positions []Position2D `json:"positions,omitempty"`
}

type LessonMusic struct {
	ElapsedTime       float32 `json:"elapsedTime"`
	Action            string  `json:"action"`
	BackgroundMusicID int64   `json:"backgroundMusicID"`
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
		lessonMaterial.BackgroundImageURL = infrastructure.GetPublicBackgroundImageURL(strconv.FormatInt(lessonMaterial.BackgroundImageID, 10))
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
