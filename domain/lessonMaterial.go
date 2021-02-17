package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type LessonMaterial struct {
	UserID            int64                `json:"userID"`
	AvatarID          int64                `json:"avatarID"`
	Duration          float32              `json:"duration"`
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
	Duration    float32    `json:"duration"`
	Position    Position3D `json:"position"`
}

type LessonGraphic struct {
	GraphicID   int64   `json:"graphicID"`
	Elapsedtime float64 `json:"elapsedtime"`
	Action      string  `json:"action"`
}

type LessonDrawing struct {
	Elapsedtime float32             `json:"elapsedtime"`
	Duration    float32             `json:"duration"`
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
	Voice       Voice    `json:"voice"`
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

/*
func GetLessonMaterialFromGCS(ctx context.Context, lessonID int64) (LessonMaterial, error) {
	lessonMaterial := new(LessonMaterial)

	bucketName := infrastructure.MaterialBucketName()
	bytes, err := infrastructure.GetObjectFromGCS(ctx, bucketName, lessonFilePath(lessonID))

	if err != nil {
		if err == storage.ErrObjectNotExist {
			return *lessonMaterial, err
		}
		return *lessonMaterial, err
	}

	if err := json.Unmarshal(bytes, lessonMaterial); err != nil {
		return *lessonMaterial, err
	}

	return *lessonMaterial, nil
}

func CreateLessonMaterialFileToGCS(ctx context.Context, lessonID int64, lessonMaterial LessonMaterial) error {
	contents, err := json.Marshal(lessonMaterial)
	if err != nil {
		return err
	}

	contentType := "application/json"
	bucketName := infrastructure.MaterialBucketName()
	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, lessonFilePath(lessonID), contentType, contents); err != nil {
		return err
	}

	return nil
}

func lessonFilePath(lessonID int64) string {
	return fmt.Sprintf("lesson/%d.json", lessonID)
}
*/
