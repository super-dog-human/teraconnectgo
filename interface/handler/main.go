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
	e.GET("/lessons", getLessons)
	e.GET("/lessons/:id", getLesson)
	e.GET("/users/:id", getUser)
	e.GET("/lesson_compressing", postLessonCompressing)
	e.POST("/lesson_compressing", postLessonCompressing)

	auth := e.Group("", Authentication())
	auth.GET("/users/me", getUserMe)
	auth.POST("/users", postUser)
	auth.PATCH("/users", patchUser)
	auth.DELETE("/users", deleteUser)
	auth.GET("/users/me/lessons", getCurrentUserLessons)
	auth.GET("/avatars", getAvatars)
	auth.POST("/avatars", postAvatars)
	auth.GET("/background_musics", getBackgroundMusics)
	auth.POST("/background_musics", postBackgroundMusic)
	auth.GET("/graphics/:id", getGraphic)
	auth.GET("/graphics", getGraphics)
	auth.POST("/graphics", postGraphics)
	auth.DELETE("/graphics/:id", deleteGraphic)
	auth.GET("/voices/:id", getVoice)
	auth.GET("/voices", getVoices)
	auth.POST("/voice", postVoice)
	auth.POST("/synthesis_voice", postSynthesisVoice)
	auth.POST("/lessons", postLesson)
	auth.PATCH("/lessons/:id", patchLesson)
	auth.GET("/lessons/:lessonID/materials/:id", getLessonMaterials)
	auth.POST("/lessons/:lessonID/materials", postLessonMaterial)
	auth.PATCH("/lessons/:lessonID/materials/:id", patchLessonMaterial)
	auth.POST("lessons/:id/thumbnail", postLessonThumbnail)

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
