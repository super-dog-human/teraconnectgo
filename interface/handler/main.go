package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Main is handling API request.
func Main(appEnv string) {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{allowOrigin(appEnv)},
	}))

	e.GET("/lessons", getLessons)
	e.GET("/lessons/:id", getLesson)
	e.GET("/avatars", getAvatars)
	e.GET("/graphics", getGraphics)

	auth := e.Group("")
	auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    domain.PublicKey(),
		SigningMethod: "RS256",
	}))

	auth.POST("/lessons", createLesson)
	auth.PATCH("/lessons/:id", updateLesson)
	auth.DELETE("/lessons/:id", destroyLesson)
	auth.GET("/lessons/:id/materials", getMaterials)
	auth.POST("/lessons/:id/materials", putMaterial)
	auth.PUT("/lessons/:id/materials", putMaterial) // same function as POST
	auth.GET("/lessons/:id/voice_texts", getVoiceTexts)
	auth.PUT("/lessons/:id/packs", updateLessonPack)
	auth.GET("/storage_objects", getStorageObjects)
	auth.POST("/storage_objects", postStorageObjects)
	auth.POST("/raw_voices", postRawVoice)

	http.Handle("/", e)
}

func allowOrigin(appEnv string) string {
	switch appEnv {
	case "production":
		return "https://authoring.teraconnect.org"
	case "staging":
		return "https://teraconnect-authoring-development-dot-teraconnect-209509.appspot.com"
	case "development":
		return "http://localhost:1234"
	default:
		return "http://localhost:1234"
	}
}