package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateSynthesisVoiceAndBlankFile creates Voice and blank files of mp3 and wav.
func CreateSynthesisVoice(request *http.Request, params *domain.CreateSynthesisVoiceParam) (domain.Voice, error) {
	ctx := request.Context()

	voice := domain.Voice{
		IsSynthesis: true,
	}

	userID, err := currentUserAccessToLesson(ctx, request, params.LessonID)
	if err != nil {
		return voice, err
	}

	voice.UserID = userID

	// ID採番のためだけにVoiceを作成する
	if err = domain.CreateVoice(ctx, params.LessonID, &voice); err != nil {
		return voice, err
	}

	mp3URL, err := domain.CreateSynthesisVoice(ctx, params, voice.ID)
	if err != nil {
		return voice, err
	}

	voice.URL = mp3URL

	return voice, nil
}
