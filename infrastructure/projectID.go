package infrastructure

// ProjectID returns Google Cloud Project ID.
func ProjectID() string {
	switch AppEnv() {
	case "production":
		return "teraconnect-209509"
	case "staging":
		return "teraconnect-staging"
	default:
		return ""
	}
}
