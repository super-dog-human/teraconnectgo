package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type AvatarErrorCode uint

const (
	AvatarNotFound AvatarErrorCode = 1
)

func (e AvatarErrorCode) Error() string {
	switch e {
	case AvatarNotFound:
		return "avatar not found"
	default:
		return "unknown avatar error"
	}
}

// Avatar is used for lesson.
type Avatar struct {
	ID       int64        `json:"id" datastore:"-"`
	Name     string       `json:"name"`
	URL      string       `json:"url"`
	Config   AvatarConfig `json:"config"`
	Version  int64        `json:"version"`
	IsPublic bool         `json:"-"`
	Created  time.Time    `json:"created"`
	Updated  time.Time    `json:"updated"`
}

type AvatarConfig struct {
	Scale        float32          `json:"scale"`
	Positions    []float32        `json:"positions"`
	InitialPoses []AvatarRotation `json:"initialPoses"`
}

type AvatarRotation struct {
	BoneName  string    `json:"boneName"`
	Rotations []float32 `json:"rotations"`
}

// GetPublicAvatarByID gets avatar by id.
func GetPublicAvatarByID(ctx context.Context, id int64) (Avatar, error) {
	avatar := new(Avatar)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *avatar, err
	}

	key := datastore.IDKey("Avatar", id, nil)
	if err := client.Get(ctx, key, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *avatar, AvatarNotFound
		} else {
			return *avatar, err
		}
	}

	avatar.ID = id
	avatar.URL = createAvatarPublicURL(id)

	return *avatar, nil
}

func GetCurrentUsersAvatarByID(ctx context.Context, id int64, userID int64) (Avatar, error) {
	avatar := new(Avatar)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *avatar, err
	}

	ancestor := datastore.IDKey("User", userID, nil)
	key := datastore.IDKey("Avatar", id, ancestor)
	if err := client.Get(ctx, key, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *avatar, AvatarNotFound
		} else {
			return *avatar, err
		}
	}

	url, err := createAvatarSignedURLs(ctx, id)
	if err != nil {
		return *avatar, err
	}

	avatar.ID = id
	avatar.URL = url

	return *avatar, nil
}

// GetCurrentUsersAvatars gets avatars belongs to user.
func GetCurrentUsersAvatars(ctx context.Context, userID int64) ([]Avatar, error) {
	var avatars []Avatar

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	ancestor := datastore.IDKey("User", userID, nil)
	query := datastore.NewQuery("Avatar").Ancestor(ancestor).Order("-Created")
	keys, err := client.GetAll(ctx, query, &avatars)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		url, err := createAvatarSignedURLs(ctx, key.ID)
		if err != nil {
			return nil, err
		}

		avatars[i].ID = key.ID
		avatars[i].URL = url
	}

	return avatars, nil
}

// GetPublicAvatars gets public avatars.
func GetPublicAvatars(ctx context.Context) ([]Avatar, error) {
	var avatars []Avatar

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return avatars, err
	}

	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true).Order("-Created")
	keys, err := client.GetAll(ctx, query, &avatars)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		avatars[i].ID = key.ID
		avatars[i].URL = createAvatarPublicURL(key.ID)
	}

	return avatars, nil
}

func createAvatarPublicURL(id int64) string {
	fileID := strconv.FormatInt(id, 10)
	return "https://storage.googleapis.com/" + infrastructure.PublicBucketName() + "/avatar/" + fileID + ".vrm"
}

func createAvatarSignedURLs(ctx context.Context, id int64) (string, error) {
	fileID := strconv.FormatInt(id, 10)
	filePath := infrastructure.StorageObjectFilePath("Avatar", fileID, "vrm")
	bucketName := infrastructure.MaterialBucketName()

	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
	if err != nil {
		return url, err
	}

	return url, nil
}

// CreateAvatar creates a new avatar belongs to user.
func CreateAvatar(ctx context.Context, avatar *Avatar, user *User) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	currentTime := time.Now()
	avatar.IsPublic = false
	avatar.Created = currentTime
	avatar.Updated = currentTime

	ancestor := datastore.IDKey("User", user.ID, nil)
	key := datastore.IncompleteKey("Avatar", ancestor)
	putKey, err := client.Put(ctx, key, avatar)
	if err != nil {
		return err
	}

	avatar.ID = putKey.ID

	return nil
}
