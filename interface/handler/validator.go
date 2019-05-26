package handler

import	"regexp"

var xidRegexp = regexp.MustCompile("[0-9a-v]{20}")

// IsValidXIDs is check valid XID string format.
func IsValidXIDs(ids []string) bool {
	regExp := xidRegexp.Copy()
	for _, id := range ids {
		if !regExp.MatchString(id) {
			return false
		}
	}

	return true
}