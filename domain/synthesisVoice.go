package domain

import (
	"context"
	"fmt"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type CreateSynthesisVoiceParam struct {
	LessonID int64  `json:"lessonID"`
	Text     string `json:"text"`
	VoiceSynthesisConfig
}

// CreateSynthesisVoice is creates new voice.
func CreateSynthesisVoice(ctx context.Context, params *CreateSynthesisVoiceParam, voice Voice) error {
	bucketName := infrastructure.PublicBucketName()
	filePath := CloudStorageVoiceFilePath(params.LessonID, voice.ID, voice.FileKey)

	if _, err := CreateSynthesizedVoice(ctx, params, bucketName, filePath); err != nil {
		return err
	}

	return nil
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

func CloudStorageVoiceFilePath(lessonID, voiceID int64, voiceFileKey string) string {
	return fmt.Sprintf("voice/%d/%d_%s.mp3", lessonID, voiceID, voiceFileKey)
}
