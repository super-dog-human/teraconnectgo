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

	e.GET("/subjects", getSubjects)
	e.GET("/categories", getCategories)
	e.GET("/background_images", getBackgroundImages)
	e.GET("/background_musics", getBackgroundMusics)
	e.GET("/lessons", getLessons)
	e.GET("/lessons/:id", getLesson)
	e.GET("/users/:id", getUser)

	auth := e.Group("", Authentication())
	auth.GET("/users/me", getUserMe)
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

	if infrastructure.AppEnv() == "development" {
		log.Fatal(http.ListenAndServeTLS(":443", "localhost.crt", "localhost.key", nil))
	} else {
		log.Fatal(e.Start(":" + port))
	}
}
