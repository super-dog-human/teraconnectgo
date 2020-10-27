package infrastructure

// ServiceAccount returns email address of google service account.
func ServiceAccount() string {
	switch AppEnv() {
	case "production":
		return ProjectID() + "@appspot.gserviceaccount.com"
	case "staging":
		return ProjectID() + "@appspot.gserviceaccount.com"
	default:
		return ""
	}
}
