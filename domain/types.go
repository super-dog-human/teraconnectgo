package domain

import "time"

type UserProviderID struct {
	ID string
}

// User is application registrated user
type User struct {
	ID         int64     `json:"id" datastore:"-"`
	ProviderID string    `json:"-"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Created    time.Time `json:"-"`
	Updated    time.Time `json:"-"`
}

// LessonReview is review status of lesson by other users.
type LessonReview struct {
	ID             int64              `json:"id" datastore:"-"`
	LessonID       int64              `json:"lessonID"`
	ReviewerUserID int64              `json:"userID"`
	Status         LessonReviewStatus `json:"status"`
	Created        time.Time          `json:"created"`
	Updated        time.Time          `json:"updated"`
}

// LessonReviewStatus is status of LessonReview.
type LessonReviewStatus uint

const (
	InReview LessonReviewStatus = 0
	Expired  LessonReviewStatus = 1
	Rejected LessonReviewStatus = 2
	Accepted LessonReviewStatus = 3
)

// Avatar is used for lesson.
type Avatar struct {
	ID              int64        `json:"id" datastore:"-"`
	UserID          int64        `json:"userID"`
	URL             string       `json:"url" datastore:"-"`
	ThumbnailURL    string       `json:"thumbnailURL" datastore:"-"`
	Name            string       `json:"name"`
	DefaultPoseKeys []AvatarPose `json:"defaultPoseKeys"`
	Version         int64        `json:"version"`
	IsPublic        bool         `json:"isPublic"`
	Created         time.Time    `json:"created"`
	Updated         time.Time    `json:"updated"`
}

// LessonAuthor is author of lesson.
type LessonAuthor struct {
	ID       int64     `json:"id" datastore:"-"`
	LessonID int64     `json:"lessonID"`
	UserID   int64     `json:"userID"`
	Role     string    `json:"role"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Graphic is used for lesson.
type Graphic struct {
	ID                int64     `json:"id" datastore:"-"`
	GraphicCategoryID int64     `json:"graphicCategoryID"`
	UserID            int64     `json:"userID"`
	FileType          string    `json:"fileType"`
	IsPublic          bool      `json:"isPublic"`
	URL               string    `json:"url" datastore:"-"`
	ThumbnailURL      string    `json:"thumbnailURL" datastore:"-"`
	Created           time.Time `json:"created"`
}

// RawVoiceText is used for lesson.
type RawVoiceText struct {
	FileID      string  `json:"fileID"`
	LessonID    int64   `json:"lessonID"`
	DurationSec float64 `json:"durationSec"`
	Text        string  `json:"text"`
	IsTexted    bool    `json:"isTexted"`
	IsConverted bool    `json:"isConverted"`
}

type LessonMaterial struct {
	DurationSec float64          `json:"durationSec"`
	Timelines   []LessonTimeline `json:"timelines"`
	PoseKey     AvatarPose       `json:"poseKey"`
	FaceKey     AvatarFace       `json:"faceKey"`
}

type LessonTimeline struct {
	TimeSec  float64         `json:"timeSec"`
	Subtitle Subtitle        `json:"subtitle"`
	Caption  Caption         `json:"caption"`
	Voice    Voice           `json:"voice"`
	Graphic  []LessonGraphic `json:"graphics"`
	Action   AvatarAction    `json:"action"`
}

type Subtitle struct {
	DurationSec float64 `json:"durationSec"`
	Body        string  `json:"body"`
}

type Caption struct {
	DurationSec     float64 `json:"durationSec"`
	Body            string  `json:"body"`
	HorizontalAlign string  `json:"horizontalAlign"`
	VerticalAlign   string  `json:"verticalAlign"`
	SizeVW          uint8   `json:"sizeVW"`
	BodyColor       string  `json:"bodyColor"`
	BorderColor     string  `json:"borderColor"`
}

type Voice struct {
	ID          string  `json:"id"`
	DurationSec float64 `json:"durationSec"`
}

type LessonGraphic struct {
	ID              string `json:"id"`
	FileType        string `json:"fileType"`
	Action          string `json:"action"`
	SizePct         uint8  `json:"sizePct"`
	HorizontalAlign string `json:"horizontalAlign"`
	VerticalAlign   string `json:"verticalAlign"`
}

type AvatarAction struct {
	Action string `json:"action"`
}

type AvatarPose struct {
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

type AvatarFace struct {
	AllAngry       AvatarWeight `json:"allAngry"`
	AllFun         AvatarWeight `json:"allFun"`
	AllJoy         AvatarWeight `json:"allJoy"`
	AllSorrow      AvatarWeight `json:"allSorrow"`
	AllSurprised   AvatarWeight `json:"allSurprised"`
	BrwAngry       AvatarWeight `json:"brwAngry"`
	BrwFun         AvatarWeight `json:"brwFun"`
	BrwJoy         AvatarWeight `json:"brwJoy"`
	BrwSorrow      AvatarWeight `json:"brwSorrow"`
	BrwSurprised   AvatarWeight `json:"brwSurprised"`
	EyeAngry       AvatarWeight `json:"eyeAngry"`
	EyeClose       AvatarWeight `json:"eyeClose"`
	EyeCloseR      AvatarWeight `json:"eyeCloseR"`
	EyeCloseL      AvatarWeight `json:"eyeCloseL"`
	EyeJoy         AvatarWeight `json:"eyeJoy"`
	EyeJoyR        AvatarWeight `json:"eyeJoyR"`
	EyeJoyL        AvatarWeight `json:"eyeJoyL"`
	EyeSorrow      AvatarWeight `json:"eyeSorrow"`
	EyeSurprised   AvatarWeight `json:"eyeSurprised"`
	EyeExtra       AvatarWeight `json:"eyeExtra"`
	MouthUp        AvatarWeight `json:"mouthUp"`
	MouthDown      AvatarWeight `json:"mouthDown"`
	MouthAngry     AvatarWeight `json:"mouthAngry"`
	MouthCorner    AvatarWeight `json:"mouthCorner"`
	MouthFun       AvatarWeight `json:"mouthFun"`
	MouthJoy       AvatarWeight `json:"mouthJoy"`
	MouthSorrow    AvatarWeight `json:"mouthSorrow"`
	MouthSurprised AvatarWeight `json:"mouthSurprised"`
	MouthA         AvatarWeight `json:"mouthA"`
	MouthI         AvatarWeight `json:"mouthI"`
	MouthU         AvatarWeight `json:"mouthU"`
	MouthE         AvatarWeight `json:"mouthE"`
	MouthO         AvatarWeight `json:"mouthO"`
	Fung1          AvatarWeight `json:"fung1"`
	Fung1Low       AvatarWeight `json:"fung1Low"`
	Fung1Up        AvatarWeight `json:"fung1Up"`
	Fung2          AvatarWeight `json:"fung2"`
	Fung2Low       AvatarWeight `json:"fung2Low"`
	Fung2Up        AvatarWeight `json:"fung2Up"`
	EyeExtraOn     AvatarWeight `json:"eyeExtraOn"`
}

type AvatarWeight struct {
	Values []float32 `json:"values"`
	Times  []float32 `json:"times"`
}
