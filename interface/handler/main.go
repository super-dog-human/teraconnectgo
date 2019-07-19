package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Main is handling API request.
func Main(appEnv string) {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{infrastructure.OriginUrl(appEnv)},
	}))

	e.GET("/lessons", getLessons)
	e.GET("/lessons/:id", getLesson)

	auth := e.Group("")
	auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    domain.PublicKey(),
		SigningMethod: "RS256",
	}))

	auth.GET("/users/me", getUserMe)
	auth.GET("/users/:id", getUser)
	auth.POST("/user", postUser)
	auth.PATCH("/user", patchUser)
	auth.DELETE("/user", deleteUser)
	auth.DELETE("/lessons/:id", deleteLesson)
	auth.GET("/avatars", getAvatars)
	auth.POST("/avatars", postAvatars)
	auth.GET("/graphics", getGraphics)
	auth.POST("/graphics", postGraphics)
	auth.GET("/authoring_lessons/:id", getAuthoringLesson)
	auth.POST("/authoring_lessons", postAuthoringLesson)
	auth.PATCH("/authoring_lessons/:id", patchAuthoringLesson)
	auth.GET("/authoring_lessons/:id/materials", getAuthoringLessonMaterials)
	auth.POST("/authoring_lessons/:id/materials", putAuthoringLessonMaterial)
	auth.PUT("/authoring_lessons/:id/materials", putAuthoringLessonMaterial) // same function as POST
	auth.GET("/lessons/:id/raw_voice_texts", getRawVoiceTexts)
	auth.PUT("/lessons/:id/packs", putLessonPack)
	auth.GET("/storage_objects", getStorageObjects)
	auth.POST("/blank_raw_voices", postBlankRawVoice)

	http.Handle("/", e)
}
