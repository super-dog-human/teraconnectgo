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
	filePath := CloudStorageVoiceFilePath(params.LessonID, voiceID)

	g.Go(func() error {
		_, err := CreateSynthesizedVoice(ctx, params, bucketName, filePath)
		return err
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

func CreateSynthesizedVoice(ctx context.Context, params *CreateSynthesisVoiceParam, bucketName, filePath string) ([]byte, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if filePath != "" {
		if err := infrastructure.CreateFileToGCS(ctx, bucketName, filePath, "audio/mpeg", resp.AudioContent); err != nil {
			return nil, err
		}
	}

	return resp.AudioContent, nil
}

func CloudStorageVoiceFilePath(lessonID, voiceID int64) string {
	return fmt.Sprintf("voice/%d/%d.mp3", lessonID, voiceID)
}

func getSignedURLOfVoiceFile(ctx context.Context, lessonID int64, voiceID int64, bucketName, filePath string, url *string) error {
	result, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
	if err != nil {
		return err
	}

	*url = result

	return nil
}
