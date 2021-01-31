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

type CreateVoiceResponse struct {
	Mp3SignedURL string `json:"mp3SignedURL"`
	WavSignedURL string `json:"wavSignedURL"`
}

// CreateVoiceAndBlankFiles creats Voice and blank files of mp3 and wav.
func CreateVoiceAndBlankFiles(request *http.Request, params *CreateVoiceParam) (CreateVoiceResponse, error) {
	ctx := request.Context()

	var response CreateVoiceResponse

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

	lessonID := strconv.FormatInt(lesson.ID, 10)
	fileID := strconv.FormatInt(voice.ID, 10)

	mp3FileRequest := domain.FileRequest{
		ID:          lessonID,
		Entity:      "voice",
		Extension:   "mp3",
		ContentType: "audio/mpeg",
	}

	mp3URL, err := domain.CreateBlankFileToGCS(ctx, fileID, "voice", mp3FileRequest)
	if err != nil {
		return response, err
	}

	wavURL, err := domain.CreateBlankFileForSpeechToTextToGCS(ctx, lessonID, fileID)
	if err != nil {
		return response, err
	}

	response.Mp3SignedURL = mp3URL
	response.WavSignedURL = wavURL

	return response, nil
}
