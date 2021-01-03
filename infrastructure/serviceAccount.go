package infrastructure

// ServiceAccountName returns email address format of google service account.
func ServiceAccountName() string {
	return ProjectID() + "@appspot.gserviceaccount.com"
}

// ServiceAccountID returns full account id.
func ServiceAccountID() string {
	return "projects/" + ProjectID() + "/serviceAccounts/" + ServiceAccountName()
}
