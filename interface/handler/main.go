package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

	e.GET("/lessons", GetLessons)
	e.GET("/lessons/:id", GetLesson)

	keyData, _ := ioutil.ReadFile("./teraconnect.pub")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)

	auth := e.Group("")
	auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: "RS256",
	}))

	auth.GET("/avatars", GetAvatars)
	auth.GET("/graphics", GetGraphics)
	auth.POST("/lessons", CreateLesson)
	auth.PATCH("/lessons/:id", UpdateLesson)
	auth.DELETE("/lessons/:id", DestroyLesson)
	auth.GET("/lessons/:id/materials", GetMaterials)
	auth.POST("/lessons/:id/materials", PutMaterial)
	auth.PUT("/lessons/:id/materials", PutMaterial) // same function as POST
	auth.GET("/lessons/:id/voice_texts", GetVoiceTexts)
	auth.PUT("/lessons/:id/packs", UpdateLessonPack)
	auth.GET("/storage_objects", GetStorageObjects)
	auth.POST("/storage_objects", PostStorageObjects)
	auth.POST("/raw_voices", PostRawVoice)

	http.Handle("/", e)
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