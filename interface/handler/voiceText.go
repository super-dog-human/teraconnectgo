package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func getVoiceTexts(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	id := c.Param("id")

	if voiceTexts, err := fetchVoiceTextsFromGCD(ctx, id); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else if len(voiceTexts) == 0 {
		log.Warningf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusNotFound, "record not found.")
	} else {
		return c.JSON(http.StatusOK, voiceTexts)
	}
}

func fetchVoiceTextsFromGCD(ctx context.Context, lessonID string) ([]domain.LessonVoiceText, error) {
	query := datastore.NewQuery("LessonVoiceText").Filter("LessonID =", lessonID).Order("FileID")

	var voiceTexts []domain.LessonVoiceText
	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return voiceTexts, err
	}

	return voiceTexts, nil
}
