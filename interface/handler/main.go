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

	auth.DELETE("/lessons/:id", destroyLesson)
	auth.GET("/avatars", getAvatars)
	auth.GET("/graphics", getGraphics)
	auth.GET("/authoring_lessons/:id", getAuthoringLesson)
	auth.POST("/authoring_lessons", createAuthoringLesson)
	auth.PATCH("/authoring_lessons/:id", updateAuthoringLesson)
	auth.GET("/authoring_lessons/:id/materials", getAuthoringLessonMaterials)
	auth.POST("/authoring_lessons/:id/materials", putAuthoringLessonMaterial)
	auth.PUT("/authoring_lessons/:id/materials", putAuthoringLessonMaterial) // same function as POST
	auth.GET("/lessons/:id/raw_voice_texts", getRawVoiceTexts)
	auth.PUT("/lessons/:id/packs", updateLessonPack)
	auth.GET("/storage_objects", getStorageObjects)
	auth.POST("/storage_objects", postStorageObjects)
	auth.POST("/raw_voices", postRawVoice)

	http.Handle("/", e)
}
