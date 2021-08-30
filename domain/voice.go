package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type VoiceErrorCode uint

const (
	VoiceNotFound VoiceErrorCode = 1
)

func (e VoiceErrorCode) Error() string {
	switch e {
	case VoiceNotFound:
		return "voice not found"
	default:
		return "unknown voice error"
	}
}

// Voice is used for lesson.
type Voice struct {
	ID          int64     `json:"id" datastore:"-"`
	UserID      int64     `json:"userID" datastore:",noindex"`
	FileKey     string    `json:"fileKey" datastore:",noindex"`
	ElapsedTime float32   `json:"elapsedTime" datastore:",noindex"`
	DurationSec float32   `json:"durationSec" datastore:",noindex"`
	Text        string    `json:"text" datastore:",noindex"`
	IsTexted    bool      `json:"isTexted" datastore:",noindex"`
	IsSynthesis bool      `json:"-" datastore:",noindex"`
	Created     time.Time `json:"created" datastore:",noindex"`
	Updated     time.Time `json:"updated" datastore:",noindex"`
}

// GetVoice is get voice entities belongs to lesson.
func GetVoice(ctx context.Context, lessonID int64, voice *Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IDKey("Voice", voice.ID, ancestor)
	if err := client.Get(ctx, key, voice); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return VoiceNotFound
		}
		return err
	}

	return nil
}

// GetVoices is get voice entities belongs to lesson.
func GetVoices(ctx context.Context, lessonID int64, voices *[]Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	query := datastore.NewQuery("Voice").Filter("IsSynthesis =", false).Ancestor(ancestor).Order("ElapsedTime")

	keys, err := client.GetAll(ctx, query, voices)
	if err != nil {
		return err
	}

	if len(*voices) == 0 {
		return VoiceNotFound
	}

	for i, key := range keys {
		(*voices)[i].ID = key.ID
	}

	// 複数のVoice取得時、署名付きURLは時間がかかりすぎるので発行しない

	return nil
}

// CreateVoice is creates new voice.
func CreateVoice(ctx context.Context, lessonID int64, voice *Voice) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	uuid, err := UUIDWithoutHypen()
	if err != nil {
		return err
	}

	voice.FileKey = uuid
	voice.Created = time.Now()

	ancestor := datastore.IDKey("Lesson", lessonID, nil)
	key := datastore.IncompleteKey("Voice", ancestor)
	putKey, err := client.Put(ctx, key, voice)
	if err != nil {
		return err
	}

	voice.ID = putKey.ID

	return nil
}
