package teraconnectgo

import (
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/SuperDogHuman/teraconnectgo/interface/handler"
	"google.golang.org/appengine"
)

// Main serve Teraconnect API
func Main(appEnv string) {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{allowOrigin(appEnv)},
	}))

	e.GET("/lessons", handler.GetLessons)
	e.GET("/lessons/:id", handler.GetLesson)

	keyData, _ := ioutil.ReadFile("teraconnect.pub")
	key, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	config := middleware.JWTConfig{
		SigningKey: key,
		SigningMethod: "RS256",
	}
	auth := e.Group("", middleware.JWT(config))
	auth.GET("/avatars", handler.GetAvatars)
	auth.GET("/graphics", handler.GetGraphics)
	auth.POST("/lessons", handler.CreateLesson)
	auth.PATCH("/lessons/:id", handler.UpdateLesson)
	auth.DELETE("/lessons/:id", handler.DestroyLesson)
	auth.GET("/lessons/:id/materials", handler.GetMaterials)
	auth.POST("/lessons/:id/materials", handler.PutMaterial)
	auth.PUT("/lessons/:id/materials", handler.PutMaterial) // same function as POST
	auth.GET("/lessons/:id/voice_texts", handler.GetVoiceTexts)
	auth.PUT("/lessons/:id/packs", handler.UpdateLessonPack)
	auth.GET("/storage_objects", handler.GetStorageObjects)
	auth.POST("/storage_objects", handler.PostStorageObjects)
	auth.POST("/raw_voices", handler.PostRawVoice)

	http.Handle("/", e)
	appengine.Main()
}

func allowOrigin(appEnv string) (string) {
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