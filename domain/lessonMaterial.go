package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonMaterial struct {
	UserID            int64                `json:"userID"`
	LessonID          int64                `json:"LessonID"`
	AvatarID          int64                `json:"avatarID"`
	DurationSec       float32              `json:"durationSec"`
	AvatarLightColor  string               `json:"avatarLightColor"`
	BackgroundImageID int64                `json:"backgroundImageID"`
	BackgroundMusicID int64                `json:"backgroundMusicID"`
	AvatarMovings     []LessonAvatarMoving `json:"avatarMovings"`
	Graphics          []LessonGraphic      `json:"graphics"`
	Drawings          []LessonDrawing      `json:"drawings"`
	Speeches          []Speech             `json:"speeches"`
	Created           time.Time            `json:"created"`
	Updated           time.Time            `json:"updated"`
	//	Status
}

type LessonAvatarMoving struct {
	Elapsedtime float32    `json:"elapsedtime"`
	DurationSec float32    `json:"durationSec"`
	Position    Position3D `json:"position"`
}

type LessonGraphic struct {
	GraphicID   int64   `json:"graphicID"`
	Elapsedtime float64 `json:"elapsedtime"`
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
	Width     int32        `json:"width"`
	Height    int32        `json:"height"`
	Color     string       `json:"color"`
	Eraser    bool         `json:"eraser"`
	LineWidth int32        `json:"lineWidth"`
	Positions []Position2D `json:"positions"`
}

type Speech struct {
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
	query := datastore.NewQuery("LessonMaterial").Ancestor(ancestor)
	var lessonMaterials []LessonMaterial
	if _, err := client.GetAll(ctx, query, &lessonMaterials); err != nil {
		return err
	}

	if len(lessonMaterials) > 0 {
		*lessonMaterial = lessonMaterials[0]
		lessonMaterial.LessonID = lessonID
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
