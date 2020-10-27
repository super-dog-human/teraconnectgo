package infrastructure

// OriginURL return API root url each current env
func OriginURL() string {
	switch AppEnv() {
	case "production":
		return "https://teraconnect.org"
	case "staging":
		return "https://teraconnect-front-dot-teraconnect-stg.an.r.appspot.com"
	default:
		return ""
	}
}
