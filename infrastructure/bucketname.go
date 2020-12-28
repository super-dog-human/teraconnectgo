package infrastructure

// MaterialBucketName is return bucket name each environments.
func MaterialBucketName() string {
	switch AppEnv() {
	case "production":
		return "teraconn_material"
	case "staging":
		return "teraconn_material_staging"
	default:
		return "teraconn_material_development"
	}
}

// RawVoiceBucketName is return bucket name each environments.
func RawVoiceBucketName() string {
	switch AppEnv() {
	case "production":
		return "teraconn_raw_voice"
	case "staging":
		return "teraconn_raw_voice_staging"
	default:
		return "teraconn_raw_voice_development"
	}
}

// VoiceForTranscriptionBucketName is return bucket name each environments.
func VoiceForTranscriptionBucketName() string {
	switch AppEnv() {
	case "production":
		return "teraconn_voice_for_transcription"
	case "staging":
		return "teraconn_voice_for_transcription_staging"
	default:
		return "teraconn_voice_for_transcription_development"
	}
}

// ThumbnailBucketName is return bucket name each environments.
func ThumbnailBucketName() string {
	switch AppEnv() {
	case "production":
		return "teraconn_thumbnail"
	case "staging":
		return "teraconn_thumbnail_staging"
	default:
		return "teraconn_thumbnail_development"
	}
}

// PublicBucketName is return public bucket name each environments.
func PublicBucketName() string {
	switch AppEnv() {
	case "production":
		return "teraconn_public"
	case "staging":
		return "teraconn_public_staging"
	default:
		return "teraconn_public_development"
	}
}
