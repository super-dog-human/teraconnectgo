package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Main is handling API request.
func Main(appEnv string) {
	infrastructure.SetAppEnv(appEnv)

	e := echo.New()
	http.Handle("/", e)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{infrastructure.OriginURL()},
	}))

	e.GET("/categories", getCategories)
	e.GET("/lessons", getLessons)
	e.GET("/lessons/:id", getLesson)

	auth := e.Group("", Authentication())
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
