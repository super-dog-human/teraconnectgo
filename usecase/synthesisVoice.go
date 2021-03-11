package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// CreateSynthesisVoiceAndBlankFile creates Voice and blank files of mp3 and wav.
func CreateSynthesisVoice(request *http.Request, params *domain.CreateSynthesisVoiceParam) (infrastructure.SignedURL, error) {
	ctx := request.Context()

	var response infrastructure.SignedURL

	userID, err := currentUserAccessToLesson(ctx, request, params.LessonID)
	if err != nil {
		return response, err
	}

	voice := domain.Voice{
		UserID: userID,
	}

	// ID採番のためだけにVoiceを作成する
	if err = domain.CreateVoice(ctx, params.LessonID, &voice); err != nil {
		return response, err
	}

	mp3URL, err := domain.CreateSynthesisVoice(ctx, params, voice.ID)
	if err != nil {
		return response, err
	}

	response.FileID = strconv.FormatInt(voice.ID, 10)
	response.SignedURL = mp3URL

	return response, nil
}
