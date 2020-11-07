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

// Category is lesson's category type.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Lesson is the lesson infomation type.
type Lesson struct {
	ID           int64     `json:"id" datastore:"-"`
	CategoryID   int64     `json:"categoryID"`
	AvatarID     int64     `json:"avatarID"`
	Avatar       Avatar    `json:"avatar" datastore:"-"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DurationSec  float64   `json:"durationSec"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	GraphicIDs   []int64   `json:"graphicIDs"`
	Graphics     []Graphic `json:"graphics" datastore:"-"`
	ViewCount    int64     `json:"viewCount"`
	Version      int64     `json:"version"`
	IsPacked     bool      `json:"isPacked"`
	IsPublic     bool      `json:"isPublic"`
	UserID       int64     `json:"userID"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

// Avatar is used for lesson.
type Avatar struct {
	ID           int64     `json:"id" datastore:"-"`
	UserID       int64     `json:"userID"`
	URL          string    `json:"url" datastore:"-"`
	ThumbnailURL string    `json:"thumbnailURL" datastore:"-"`
	Name         string    `json:"name"`
	Version      int64     `json:"version"`
	IsPublic     bool      `json:"isPublic"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
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
	URL               string    `json:"url" datastore:"-"`
	ThumbnailURL      string    `json:"thumbnailURL" datastore:"-"`
	GraphicCategoryID int64     `json:"graphicCategoryID"`
	UserID            int64     `json:"userID"`
	FileType          string    `json:"fileType"`
	IsPublic          bool      `json:"isPublic"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
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
