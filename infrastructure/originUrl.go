package infrastructure

// OriginURL return API root url each current env
func OriginURL() string {
	switch AppEnv() {
	case "production":
		return "https://teraconnect.org"
	case "staging":
		return "https://staging.teraconnect.org"
	default:
		return "https://dev.teraconnect.org:3000"
	}
}
