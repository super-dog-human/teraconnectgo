package usecase

import (
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/jinzhu/copier"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type UserErrorCode uint

const (
	UserNotAvailable  UserErrorCode = 1
	AlreadyUserExists UserErrorCode = 2
	UserNotFound      UserErrorCode = 3
)

func (e UserErrorCode) Error() string {
	switch e {
	case UserNotAvailable:
		return "user not available"
	case AlreadyUserExists:
		return "user is already created"
	case UserNotFound:
		return "user not found"
	default:
		return "unknown error"
	}
}

// NewUserParamsは、Userの新規作成時、リクエストボディをbindするために使用されます。
type NewUserParams struct {
	Name    string `json:"name" datastore:",noindex"`
	Profile string `json:"profile" datastore:",noindex"`
	Email   string `json:"email" datastore:",noindex"`
}

// GetCurrentUser for fetch current user account
func GetCurrentUser(request *http.Request) (domain.User, error) {
	var currentUser domain.User
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return currentUser, err
	}

	return currentUser, nil
}

// GetUserはidからユーザーを取得して返します。
func GetUser(request *http.Request, id int64) (domain.User, error) {
	ctx := request.Context()
	var user domain.User

	user, err := domain.GetUserByID(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUsersは複数のユーザーを取得して返します。cursorStrがあればページネーションに使用し、なければ1件目からユーザーを返します。
func GetUsers(request *http.Request, cursorStr string) ([]domain.User, string, error) {
	ctx := request.Context()

	users, nextCursorStr, err := domain.GetUsers(ctx, cursorStr)
	if err != nil {
		return nil, "", err
	}

	if len(users) == 0 {
		return nil, "", UserNotFound
	}

	return users, nextCursorStr, nil
}

// CreateUserは新規ユーザーを作成します。
func CreateUser(request *http.Request, newUser *NewUserParams) error {
	ctx := request.Context()

	var user domain.User
	copier.Copy(&user, &newUser)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	providerID, err := domain.ProviderID(request)
	if err != nil {
		return err
	}
	user.ProviderID = providerID

	backgroundImages, err := domain.GetAllBackgroundImages(ctx)
	if err != nil {
		return err
	}

	imageIndex := time.Now().UnixNano() / 1000 % int64(len(backgroundImages))
	backgroundImage := backgroundImages[imageIndex]
	if backgroundImage.Name == "学習机" { // この画像はヘッダー画像に向かないので使用しない
		if imageIndex == 0 {
			imageIndex += 1
		} else {
			imageIndex -= 1
		}
		backgroundImage = backgroundImages[imageIndex]
	}
	user.BackgroundImageID = backgroundImage.ID

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err = domain.ReserveUserProviderIDInTransaction(tx, providerID)
		if err == domain.AlreadyProviderIDExists {
			return AlreadyUserExists
		}
		if err != nil {
			return err
		}

		if _, err = domain.CreateUserInTransaction(tx, &user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(request *http.Request, params *map[string]interface{}) (domain.User, error) {
	ctx := request.Context()

	user, err := domain.GetCurrentUser(request)
	if err != nil {
		return user, UserNotAvailable
	}

	targetFields := []string{"Name", "Profile", "Email"}
	if err = domain.UpdateUserByJson(ctx, &user, params, &targetFields); err != nil {
		return user, err
	}

	return user, nil
}

func UnsubscribeCurrentUser(request *http.Request) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	if err = domain.DeleteUser(ctx, currentUser.ID); err != nil {
		return err
	}

	return nil
}
