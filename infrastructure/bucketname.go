package infrastructure

import (
	"context"
)

// MaterialBucketName is return bucket name each environments.
func MaterialBucketName(ctx context.Context) string {
	switch AppEnv() {
	case "production":
		return "teraconn_material"
	case "staging":
		return "teraconn_material_staging"
	default:
		return ""
	}
}

// RawVoiceBucketName is return bucket name each environments.
func RawVoiceBucketName(ctx context.Context) string {
	switch AppEnv() {
	case "production":
		return "teraconn_raw_voice"
	case "staging":
		return "teraconn_raw_voice_staging"
	default:
		return ""
	}
}

// VoiceForTranscriptionBucketName is return bucket name each environments.
func VoiceForTranscriptionBucketName(ctx context.Context) string {
	switch AppEnv() {
	case "production":
		return "teraconn_voice_for_transcription"
	case "staging":
		return "teraconn_voice_for_transcription_staging"
	default:
		return ""
	}
}
