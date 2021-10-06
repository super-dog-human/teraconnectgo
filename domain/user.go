package domain

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
	"google.golang.org/api/iterator"
)

type UserProviderID struct {
	ID int64 `datastore:"-"`
}

// User is application registrated user
type User struct {
	ID                      int64     `json:"id" datastore:"-"`
	ProviderID              string    `json:"-"`
	BackgroundImageID       int64     `json:"-" datastore:",noindex"`
	BackgroundImageURL      string    `json:"backgroundImageURL" datastore:"-"`
	IntroductionID          int64     `json:"introductionID" datastore:",noindex"`
	IsPublishedIntroduction bool      `json:"isPublishedIntroduction" datastore:",noindex"`
	Name                    string    `json:"name" datastore:",noindex"`
	Profile                 string    `json:"profile" datastore:",noindex"`
	Email                   string    `json:"email,omitempty" datastore:",noindex"`
	Created                 time.Time `json:"-"`
	Updated                 time.Time `json:"-" datastore:",noindex"`
}

// UserErrorCode is user error code.
type UserErrorCode uint

const (
	// AlreadyProviderIDExists is exists privider-id of user
	AlreadyProviderIDExists UserErrorCode = 1
)

func (e UserErrorCode) Error() string {
	switch e {
	case AlreadyProviderIDExists:
		return "provider id is already existed"
	default:
		return "unknown error"
	}
}

// GetCurrentUser returns user from valid token.
func GetCurrentUser(request *http.Request) (User, error) {
	user := new(User) // for return blank user when error

	providerID, err := ProviderID(request)
	if err != nil {
		return *user, err
	}

	var users []User
	ctx := request.Context()
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, FailedDatastoreInitialize
	}

	query := datastore.NewQuery("User").Filter("ProviderID =", providerID).Limit(1)
	keys, err := client.GetAll(ctx, query, &users)
	if err != nil {
		return *user, FailedGettingUser
	}

	if len(users) == 0 {
		return *user, UserNotFound
	}

	user = &users[0]
	user.ID = keys[0].ID

	return *user, nil
}

// GetUserByIDはidからユーザーを取得して返します。Emailは必ず空文字列になり、json出力時はフィールドごとなくなります。
func GetUserByID(ctx context.Context, id int64) (User, error) {
	user := new(User)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, err
	}

	key := datastore.IDKey("User", id, nil)
	if err := client.Get(ctx, key, user); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *user, UserNotFound
		}
		return *user, err
	}
	user.ID = id
	user.BackgroundImageURL = infrastructure.GetPublicBackgroundImageURL(strconv.FormatInt(user.BackgroundImageID, 10))
	user.Email = "" // メールアドレスは返さない

	return *user, nil
}

func GetUsers(ctx context.Context, cursorStr string) ([]User, string, error) {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, "", err
	}

	const userPageSize = 20
	query := datastore.NewQuery("User").Order("-Created").Limit(userPageSize)

	if cursorStr != "" {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return nil, "", err
		}
		query = query.Start(cursor)
	}

	var users []User
	it := client.Run(ctx, query)
	for {
		var user User
		key, err := it.Next(&user)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, "", err
		}
		user.ID = key.ID
		user.Email = "" // メールアドレスは返さない
		users = append(users, user)
	}

	nextCursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}

	return users, nextCursor.String(), nil
}

// ReserveUserProviderIDInTransaction creates user's ProviderID for exclusion control.
func ReserveUserProviderIDInTransaction(tx *datastore.Transaction, providerID string) error {
	key := datastore.NameKey("UserProviderID", providerID, nil)
	userProviderID := new(UserProviderID)

	err := tx.Get(key, userProviderID)
	if err == nil {
		return AlreadyProviderIDExists
	}
	if err != datastore.ErrNoSuchEntity {
		return err
	}

	// Put only when ErrNoSuchEntity
	_, err = tx.Put(key, userProviderID)
	return err
}

// CreateUserInTransaction creates new user.
func CreateUserInTransaction(tx *datastore.Transaction, user *User) (*datastore.PendingKey, error) {
	key := datastore.IncompleteKey("User", nil)

	currentTime := time.Now()
	user.Created = currentTime
	user.Updated = currentTime

	pendingKey, err := tx.Put(key, user)
	if err != nil {
		return nil, err
	}

	return pendingKey, nil
}

// UpdateUserByJsonは、json構造のinterfaceを受け取り、Userを更新します。
func UpdateUserByJson(ctx context.Context, user *User, jsonBody *map[string]interface{}, targetFields *[]string) error {
	MergeJsonToStruct(jsonBody, user, targetFields)

	if err := UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

// UpdateUserは、受け取ったUserでエンティティを更新します。
func UpdateUser(ctx context.Context, user *User) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.IDKey("User", user.ID, nil)
	user.Updated = time.Now()

	if _, err := client.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes user.
func DeleteUser(ctx context.Context, id int64) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.IDKey("User", id, nil)
	if err := client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
