package domain

import (
	"context"
	"fmt"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
	"golang.org/x/sync/errgroup"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type CreateSynthesisVoiceParam struct {
	LessonID int64  `json:"lessonID"`
	Text     string `json:"text"`
	VoiceSynthesisConfig
}

// CreateSynthesisVoice is creates new voice.
func CreateSynthesisVoice(ctx context.Context, params *CreateSynthesisVoiceParam, voiceID int64) (string, error) {
	g, ctx := errgroup.WithContext(ctx)

	bucketName := infrastructure.MaterialBucketName()
	filePath := fmt.Sprintf("voice/%d/%d.mp3", params.LessonID, voiceID)

	g.Go(func() error {
		return createSynthesizedVoice(ctx, params, bucketName, filePath)
	})

	var url string
	g.Go(func() error {
		return getSignedURLOfVoiceFile(ctx, params.LessonID, voiceID, bucketName, filePath, &url)
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	return url, nil
}

func createSynthesizedVoice(ctx context.Context, params *CreateSynthesisVoiceParam, bucketName, filePath string) error {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: params.Text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: params.LanguageCode,
			Name:         params.Name,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  params.SpeakingRate,
			Pitch:         params.Pitch,
			VolumeGainDb:  params.VolumeGainDb,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return err
	}

	if err := infrastructure.CreateFileToGCS(ctx, bucketName, filePath, "audio/mpeg", resp.AudioContent); err != nil {
		return err
	}

	return nil
}

func getSignedURLOfVoiceFile(ctx context.Context, lessonID int64, voiceID int64, bucketName, filePath string, url *string) error {
	result, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
	if err != nil {
		return err
	}

	*url = result

	return nil
}
