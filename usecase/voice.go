package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
)

type CreateVoiceParam struct {
	LessonID    int64   `json:"lessonID"`
	Speeched    float64 `json:"speeched"`
	DurationSec float64 `json:"durationSec"`
}

// CreateVoiceAndBlankFile creats Voice and blank files of mp3 and wav.
func CreateVoiceAndBlankFile(request *http.Request, params *CreateVoiceParam) (domain.SignedURL, error) {
	ctx := request.Context()

	var response domain.SignedURL

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return response, err
	}

	lesson, err := domain.GetLessonByID(ctx, params.LessonID)
	if err != nil {
		return response, err
	}

	if lesson.UserID != currentUser.ID {
		return response, LessonNotAvailable
	}

	voice := domain.Voice{
		LessonID:    lesson.ID,
		Speeched:    params.Speeched,
		DurationSec: params.DurationSec,
	}

	if err = domain.CreateVoice(ctx, &currentUser, &voice); err != nil {
		return response, err
	}

	voiceID := strconv.FormatInt(voice.ID, 10)

	mp3FileRequest := domain.FileRequest{
		ID:          voiceID,
		Entity:      "voice",
		Extension:   "mp3",
		ContentType: "audio/mpeg",
	}

	mp3URL, err := domain.CreateBlankFileToGCS(ctx, voiceID, "voice", mp3FileRequest)
	if err != nil {
		return response, err
	}

	response.FileID = voiceID
	response.SignedURL = mp3URL

	return response, nil
}
