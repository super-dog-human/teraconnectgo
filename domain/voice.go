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
	UserID      int64     `json:"userID"`
	LessonID    int64     `json:"lessonID"`
	Elapsedtime float32   `json:"elapsedtime"`
	DurationSec float32   `json:"durationSec"`
	Text        string    `json:"text"`
	IsTexted    bool      `json:"isTexted"`
	URL         string    `json:"url" datastore:"-"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// CreateVoice is creats new voice.
func CreateVoice(ctx context.Context, voice *Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	voice.Created = time.Now()

	key := datastore.IncompleteKey("Voice", nil)
	putKey, err := client.Put(ctx, key, voice)
	if err != nil {
		return err
	}

	voice.ID = putKey.ID

	return nil
}
