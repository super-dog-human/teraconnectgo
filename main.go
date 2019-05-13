package teraconnectgo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/SuperDogHuman/teraconnectgo/interface"
	"google.golang.org/appengine"
)

// Main serve Teraconnect API
func Main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			//"https://authoring.teraconnect.org",
			//"https://teraconnect-authoring-development-dot-teraconnect-209509.appspot.com",
			"http://localhost:1234",
		},
	}))

	e.GET("/lessons", interface.GetLessons)
	e.GET("/lessons/:id", interface.GetLesson)

	auth := e.Group("", middleware.JWT([]byte("secret")))
	auth.GET("/avatars", interface.GetAvatars)
	auth.GET("/graphics", interface.GetGraphics)
	auth.POST("/lessons", interface.CreateLesson)
	auth.PATCH("/lessons/:id", interface.UpdateLesson)
	auth.DELETE("/lessons/:id", interface.DestroyLesson)
	auth.GET("/lessons/:id/materials", interface.GetLessonMaterials)
	auth.POST("/lessons/:id/materials", interface.PutLessonMaterial)
	auth.PUT("/lessons/:id/materials", interface.PutLessonMaterial) // same function as POST
	auth.GET("/lessons/:id/voice_texts", interface.GetVoiceTexts)
	auth.PUT("/lessons/:id/packs", interface.UpdateLessonPack)
	auth.GET("/storage_objects", interface.GetStorageObjects)
	auth.POST("/storage_objects", interface.PostStorageObjects)
	auth.POST("/raw_voices", interface.PostRawVoices)

	http.Handle("/", e)
	appengine.Main()
}