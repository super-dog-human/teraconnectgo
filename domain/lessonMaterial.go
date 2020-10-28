package domain

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/storage"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetLessonMaterialFromGCS(ctx context.Context, lessonID string) (LessonMaterial, error) {
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

func CreateLessonMaterialFileToGCS(ctx context.Context, lessonID string, lessonMaterial LessonMaterial) error {
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

func lessonFilePath(lessonID string) string {
	return "lesson/" + lessonID + ".json"
}

type LessonMaterial struct {
	DurationSec float64          `json:"durationSec"`
	Timelines   []LessonTimeline `json:"timelines"`
	PoseKey     LessonAvatarPose `json:"poseKey"`
	FaceKey     LessonAvatarFace `json:"faceKey"`
}

type LessonTimeline struct {
	TimeSec float64                 `json:"timeSec"`
	Text    LessonMaterialText      `json:"text"`
	Voice   LessonMaterialVoice     `json:"voice"`
	Graphic []LessonMaterialGraphic `json:"graphics"`
	Action  LessonAvatarAction      `json:"action"`
}

type LessonMaterialText struct {
	DurationSec     float64 `json:"durationSec"`
	Body            string  `json:"body"`
	HorizontalAlign string  `json:"horizontalAlign"`
	VerticalAlign   string  `json:"verticalAlign"`
	SizeVW          uint8   `json:"sizeVW"`
	BodyColor       string  `json:"bodyColor"`
	BorderColor     string  `json:"borderColor"`
}

type LessonMaterialVoice struct {
	ID          string  `json:"id"`
	DurationSec float64 `json:"durationSec"`
}

type LessonMaterialGraphic struct {
	ID              string `json:"id"`
	FileType        string `json:"fileType"`
	Action          string `json:"action"`
	SizePct         uint8  `json:"sizePct"`
	HorizontalAlign string `json:"horizontalAlign"`
	VerticalAlign   string `json:"verticalAlign"`
}

type LessonAvatarAction struct {
	Action string `json:"action"`
}

type LessonAvatarPose struct {
	LeftHands      []LessonRotation `json:"leftHands"`
	RightHands     []LessonRotation `json:"rightHands"`
	LeftElbows     []LessonRotation `json:"leftElbows"`
	RightElbows    []LessonRotation `json:"rightElbows"`
	LeftShoulders  []LessonRotation `json:"leftShoulders"`
	RightShoulders []LessonRotation `json:"rightShoulders"`
	Necks          []LessonRotation `json:"necks"`
	CoreBodies     []LessonPosition `json:"coreBodies"`
}

type LessonRotation struct {
	Rot  []float32 `json:"rot"`
	Time float32   `json:"time"`
}

type LessonPosition struct {
	Pos  []float32 `json:"pos"`
	Time float32   `json:"time"`
}

type LessonAvatarFace struct {
	AllAngry       LessonWeight `json:"allAngry"`
	AllFun         LessonWeight `json:"allFun"`
	AllJoy         LessonWeight `json:"allJoy"`
	AllSorrow      LessonWeight `json:"allSorrow"`
	AllSurprised   LessonWeight `json:"allSurprised"`
	BrwAngry       LessonWeight `json:"brwAngry"`
	BrwFun         LessonWeight `json:"brwFun"`
	BrwJoy         LessonWeight `json:"brwJoy"`
	BrwSorrow      LessonWeight `json:"brwSorrow"`
	BrwSurprised   LessonWeight `json:"brwSurprised"`
	EyeAngry       LessonWeight `json:"eyeAngry"`
	EyeClose       LessonWeight `json:"eyeClose"`
	EyeCloseR      LessonWeight `json:"eyeCloseR"`
	EyeCloseL      LessonWeight `json:"eyeCloseL"`
	EyeJoy         LessonWeight `json:"eyeJoy"`
	EyeJoyR        LessonWeight `json:"eyeJoyR"`
	EyeJoyL        LessonWeight `json:"eyeJoyL"`
	EyeSorrow      LessonWeight `json:"eyeSorrow"`
	EyeSurprised   LessonWeight `json:"eyeSurprised"`
	EyeExtra       LessonWeight `json:"eyeExtra"`
	MouthUp        LessonWeight `json:"mouthUp"`
	MouthDown      LessonWeight `json:"mouthDown"`
	MouthAngry     LessonWeight `json:"mouthAngry"`
	MouthCorner    LessonWeight `json:"mouthCorner"`
	MouthFun       LessonWeight `json:"mouthFun"`
	MouthJoy       LessonWeight `json:"mouthJoy"`
	MouthSorrow    LessonWeight `json:"mouthSorrow"`
	MouthSurprised LessonWeight `json:"mouthSurprised"`
	MouthA         LessonWeight `json:"mouthA"`
	MouthI         LessonWeight `json:"mouthI"`
	MouthU         LessonWeight `json:"mouthU"`
	MouthE         LessonWeight `json:"mouthE"`
	MouthO         LessonWeight `json:"mouthO"`
	Fung1          LessonWeight `json:"fung1"`
	Fung1Low       LessonWeight `json:"fung1Low"`
	Fung1Up        LessonWeight `json:"fung1Up"`
	Fung2          LessonWeight `json:"fung2"`
	Fung2Low       LessonWeight `json:"fung2Low"`
	Fung2Up        LessonWeight `json:"fung2Up"`
	EyeExtraOn     LessonWeight `json:"eyeExtraOn"`
}

type LessonWeight struct {
	Values []float32 `json:"values"`
	Times  []float32 `json:"times"`
}
