package infrastructure

import (
	"context"

	"google.golang.org/appengine"
)

// MaterialBucketName is return bucket name each environments.
func MaterialBucketName(ctx context.Context) string {
	bucketName := "teraconn_material"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return bucketName
	}

	return bucketName + "_development"
}

// RawVoiceBucketName is return bucket name each environments.
func RawVoiceBucketName(ctx context.Context) string {
	bucketName := "teraconn_raw_voice"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return bucketName
	}

	return bucketName + "_development"
}

// VoiceForTranscriptionBucketName is return bucket name each environments.
func VoiceForTranscriptionBucketName(ctx context.Context) string {
	bucketName := "teraconn_voice_for_transcription"
	module := appengine.ModuleName(ctx)

	if module == "teraconnect-api" {
		return bucketName
	}

	return bucketName + "_development"
}
