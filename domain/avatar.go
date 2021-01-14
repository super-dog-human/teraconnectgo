package domain

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

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

type AvatarPose struct {
	LeftHands      []LessonRotation `json:"leftHands"`
	RightHands     []LessonRotation `json:"rightHands"`
	LeftElbows     []LessonRotation `json:"leftElbows"`
	RightElbows    []LessonRotation `json:"rightElbows"`
	LeftShoulders  []LessonRotation `json:"leftShoulders"`
	RightShoulders []LessonRotation `json:"rightShoulders"`
	Necks          []LessonRotation `json:"necks"`
	CoreBodies     []LessonPosition `json:"coreBodies"`
}

type LessonRotation struct {
	Rot  []float32 `json:"rot"`
	Time float32   `json:"time"`
}

type LessonPosition struct {
	Pos  []float32 `json:"pos"`
	Time float32   `json:"time"`
}

// GetAvatarByIDs gets avatar by id.
func GetAvatarByIDs(ctx context.Context, id int64) (Avatar, error) {
	avatar := new(Avatar)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *avatar, err
	}

	key := datastore.IDKey("Avatar", id, nil)
	if err := client.Get(ctx, key, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *avatar, err
		}
		return *avatar, err
	}

	avatar.ID = id

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
		avatars[i].URL = createAvatarPublicURLs(key.ID)
	}

	return avatars, nil
}

func createAvatarPublicURLs(id int64) string {
	fileID := strconv.FormatInt(id, 10)
	return "https://storage.googleapis.com/" + infrastructure.MaterialBucketName() + "/avatar/" + fileID + ".vrm"
}

func createAvatarSignedURLs(ctx context.Context, id int64) (string, error) {
	fileID := strconv.FormatInt(id, 10)
	filePath := storageObjectFilePath("Avatar", fileID, "vrm")
	bucketName := infrastructure.MaterialBucketName()

	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
	if err != nil {
		return url, err
	}

	return url, nil
}

// CreateAvatar creats a new avatar belongs to user.
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

func DeleteAvatarInTransaction(tx *datastore.Transaction, id int64) error {
	key := datastore.IDKey("Avatar", id, nil)
	if err := tx.Delete(key); err != nil {
		return err
	}

	return nil
}
