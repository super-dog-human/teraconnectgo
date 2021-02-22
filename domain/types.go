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

type Position2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Position3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
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
