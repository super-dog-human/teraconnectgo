package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Voice is used for lesson.
type Voice struct {
	ID          int64     `json:"id" datastore:"-"`
	LessonID    int64     `json:"lessonID"`
	Speeched    float64   `json:"speeched"`
	DurationSec float64   `json:"durationSec"`
	Text        string    `json:"text"`
	IsTexted    bool      `json:"isTexted"`
	URL         string    `json:"url" datastore:"-"`
	Created     time.Time `json:"created"`
}

// CreateVoice is creats new voice.
func CreateVoice(ctx context.Context, user *User, voice *Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	voice.Created = time.Now()

	parentKey := datastore.IDKey("User", user.ID, nil)
	key := datastore.IncompleteKey("Voice", parentKey)
	putKey, err := client.Put(ctx, key, voice)
	if err != nil {
		return err
	}

	voice.ID = putKey.ID

	return nil
}
