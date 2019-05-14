package interface

import (
	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	//"google.golang.org/appengine/memcache"
	"cloudHelper"
	"encoding/json"
	"net/http"
	"utility"
)

// GetMaterials is get material of the lesson function.
func GetMaterials(c echo.Context) error {
	// increment view cont in memorycache
	// https://cloud.google.com/appengine/docs/standard/go/memcache/reference
	// https://cloud.google.com/appengine/docs/standard/go/memcache/using?hl=ja

	lessonID := c.Param("id")
	ctx := appengine.NewContext(c.Request())

	ids := []string{lessonID}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	filePath := "lesson/" + lessonID + ".json"
	bucketName := utility.MaterialBucketName(ctx)

	bytes, err := cloudHelper.GetObjectFromGCS(ctx, bucketName, filePath)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			log.Warningf(ctx, "%+v\n", errors.WithStack(err))
			return c.JSON(http.StatusNotFound, err.Error())
		}
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	lessonMaterial := new(LessonMaterial)
	if err := json.Unmarshal(bytes, lessonMaterial); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lessonMaterial)
}

// PutMaterial is put material of the lesson function.
func PutMaterial(c echo.Context) error {
	lessonID := c.Param("id")
	ctx := appengine.NewContext(c.Request())

	ids := []string{lessonID}
	if !utility.IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		log.Warningf(ctx, errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial := new(LessonMaterial)

	if err := c.Bind(lessonMaterial); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	contents, err := json.Marshal(lessonMaterial)
	if err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	filePath := "lesson/" + lessonID + ".json"
	contentType := "application/json"
	bucketName := utility.MaterialBucketName(ctx)

	if err := cloudHelper.CreateObjectToGCS(ctx, bucketName, filePath, contentType, contents); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeed")
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
