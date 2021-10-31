package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type CreateVoiceParam struct {
	LessonID    int64   `json:"lessonID"`
	ElapsedTime float32 `json:"elapsedTime"`
	DurationSec float32 `json:"durationSec"`
}

func GetVoices(request *http.Request, lessonID int64) ([]domain.Voice, error) {
	ctx := request.Context()

	var voices []domain.Voice

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return nil, err
	}

	if err := domain.GetVoices(ctx, lessonID, &voices); err != nil {
		return nil, err
	}

	return voices, nil
}

// CreateVoiceAndBlankFile creates Voice and blank mp3 file.
func CreateVoiceAndBlankFile(request *http.Request, params *CreateVoiceParam) (domain.Voice, string, error) {
	ctx := request.Context()

	var voice domain.Voice

	userID, err := currentUserAccessToLesson(ctx, request, params.LessonID)
	if err != nil {
		return voice, "", err
	}

	voice.UserID = userID
	voice.LessonID = params.LessonID
	voice.ElapsedTime = params.ElapsedTime
	voice.DurationSec = params.DurationSec

	if err = domain.CreateVoice(ctx, &voice); err != nil {
		return voice, "", err
	}

	lessonID := strconv.FormatInt(params.LessonID, 10)
	fileName := strconv.FormatInt(voice.ID, 10) + "_" + voice.FileKey

	mp3FileRequest := infrastructure.FileRequest{
		ID:          fileName,
		Entity:      "voice",
		Extension:   "mp3",
		ContentType: "audio/mpeg",
	}

	filePath := lessonID + "/" + fileName
	mp3URL, err := infrastructure.CreateBlankFileToPublicGCS(ctx, filePath, "voice", mp3FileRequest)
	if err != nil {
		return voice, mp3URL, err
	}

	return voice, mp3URL, nil
}
