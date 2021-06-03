package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type CreateBackgroundMusicParam struct {
	Name string `json:"name"`
}

// GetBackgroundMusics returns music URLs in Cloud Datastore.
func GetBackgroundMusics(request *http.Request) ([]domain.BackgroundMusic, error) {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	var musics []domain.BackgroundMusic

	usersMusics, err := domain.GetCurrentUsersBackgroundMusics(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	musics = append(musics, usersMusics...)

	publicMusics, err := domain.GetPublicBackgroundMusics(ctx)
	if err != nil {
		return nil, err
	}
	musics = append(musics, publicMusics...)

	return musics, nil
}

func CreateBackgroundMusicAndBlankFile(request *http.Request, param *CreateBackgroundMusicParam) (infrastructure.SignedURL, error) {
	ctx := request.Context()

	var signedURL infrastructure.SignedURL

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURL, err
	}

	backgroundMusic := new(domain.BackgroundMusic)
	backgroundMusic.Name = param.Name
	backgroundMusic.IsPublic = false

	if err = domain.CreateBackgroundMusic(ctx, currentUser.ID, backgroundMusic); err != nil {
		return signedURL, err
	}

	fileID := strconv.FormatInt(backgroundMusic.ID, 10)
	bgmID := strconv.FormatInt(backgroundMusic.ID, 10)
	mp3FileRequest := infrastructure.FileRequest{
		ID:          bgmID,
		Entity:      "bgm",
		Extension:   "mp3",
		ContentType: "audio/mpeg",
	}

	url, err := infrastructure.CreateBlankFileToGCS(ctx, fileID, "bgm", mp3FileRequest)
	if err != nil {
		return signedURL, err
	}
	signedURL = infrastructure.SignedURL{FileID: fileID, SignedURL: url}

	return signedURL, nil
}
