package infrastructure

// ServiceAccount returns email address of google service account.
func ServiceAccount() string {
	return ProjectID() + "@appspot.gserviceaccount.com"
}
