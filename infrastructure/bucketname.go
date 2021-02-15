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
		return "teraconn_public_staging_2"
	default:
		return "teraconn_public_development"
	}
}
