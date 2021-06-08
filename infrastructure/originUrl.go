package infrastructure

// OriginURL return API root url each current env
func OriginURL() string {
	switch AppEnv() {
	case "production":
		return "https://teraconnect.org"
	case "staging":
		return "https://teraconnect-staging.an.r.appspot.com"
	default:
		return "https://dev.teraconnect.org:3000"
	}
}
