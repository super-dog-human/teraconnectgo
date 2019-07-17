package infrastructure

// OriginUrl return API root url each current env
func OriginUrl(appEnv string) string {
	switch appEnv {
	case "production":
		return "https://authoring.teraconnect.org"
	case "staging":
		return "https://teraconnect-authoring-development-dot-teraconnect-209509.appspot.com"
	case "development":
		return "http://localhost:1234"
	default:
		return "http://localhost:1234"
	}
}
