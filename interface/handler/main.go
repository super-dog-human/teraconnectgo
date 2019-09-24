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

	e.GET("/categories", getCategories)
	e.GET("/lessons", getLessons)
	e.GET("/lessons/:id", getLesson)

	auth := e.Group("")
	auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    domain.PublicKey(),
		SigningMethod: "RS256",
	}))

	auth.GET("/users/me", getUserMe)
	auth.GET("/users/:id", getUser) // publicでもいいかも
	auth.POST("/users", postUser)
	auth.PATCH("/users", patchUser)
	auth.DELETE("/users", deleteUser)
	auth.GET("/avatars", getAvatars)
	auth.POST("/avatars", postAvatars)
	auth.GET("/graphics", getGraphics)
	auth.POST("/graphics", postGraphics)
	auth.POST("/lessons", postLesson)
	auth.PATCH("/lessons/:id", patchLesson)
	auth.DELETE("/lessons/:id", deleteLesson)
	auth.GET("/lessons/:id/materials", getLessonMaterials)
	auth.POST("/lessons/:id/materials", postLessonMaterial)
	auth.PUT("/lessons/:id/materials", putLessonMaterial)
	auth.GET("/lessons/:id/raw_voice_texts", getRawVoiceTexts)
	auth.PUT("/lessons/:id/packs", putLessonPack)
	auth.GET("/storage_objects", getStorageObjects)
	auth.POST("/blank_raw_voices", postBlankRawVoice)

	http.Handle("/", e)
}
