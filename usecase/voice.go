package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
)

type CreateVoiceParam struct {
	LessonID    int64   `json:"lessonID"`
	Elapsedtime float32 `json:"elapsedtime"`
	DurationSec float32 `json:"durationSec"`
}

func GetVoices(request *http.Request, lessonID int64) ([]domain.Voice, error) {
	ctx := request.Context()

	var voices []domain.Voice

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return voices, err
	}

	if err := domain.GetVoices(ctx, lessonID, &voices); err != nil {
		return voices, err
	}

	return voices, nil
}

// CreateVoiceAndBlankFile creats Voice and blank files of mp3 and wav.
func CreateVoiceAndBlankFile(request *http.Request, params *CreateVoiceParam) (domain.SignedURL, error) {
	ctx := request.Context()

	var response domain.SignedURL

	userID, err := currentUserAccessToLesson(ctx, request, params.LessonID)
	if err != nil {
		return response, err
	}

	voice := domain.Voice{
		UserID:      userID,
		Elapsedtime: params.Elapsedtime,
		DurationSec: params.DurationSec,
	}

	if err = domain.CreateVoice(ctx, params.LessonID, &voice); err != nil {
		return response, err
	}

	lessonID := strconv.FormatInt(params.LessonID, 10)
	voiceID := strconv.FormatInt(voice.ID, 10)

	mp3FileRequest := domain.FileRequest{
		ID:          voiceID,
		Entity:      "voice",
		Extension:   "mp3",
		ContentType: "audio/mpeg",
	}

	filePath := lessonID + "/" + voiceID
	mp3URL, err := domain.CreateBlankFileToGCS(ctx, filePath, "voice", mp3FileRequest)
	if err != nil {
		return response, err
	}

	response.FileID = voiceID
	response.SignedURL = mp3URL

	return response, nil
}
