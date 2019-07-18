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

func getRawVoiceTexts(c echo.Context) error {
	ctx := appengine.NewContext(c.Request())
	id := c.Param("id")

	if voiceTexts, err := fetchVoiceTexts(ctx, id); err != nil {
		log.Errorf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else if len(voiceTexts) == 0 {
		log.Warningf(ctx, "%+v\n", errors.WithStack(err))
		return c.JSON(http.StatusNotFound, "record not found.")
	} else {
		return c.JSON(http.StatusOK, voiceTexts)
	}
}

func fetchVoiceTexts(ctx context.Context, lessonID string) ([]domain.RawVoiceText, error) {
	query := datastore.NewQuery("RawVoiceText").Filter("LessonID =", lessonID).Order("FileID")

	var voiceTexts []domain.RawVoiceText
	if _, err := query.GetAll(ctx, &voiceTexts); err != nil {
		return voiceTexts, err
	}

	return voiceTexts, nil
}
