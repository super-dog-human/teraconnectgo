package domain

import (
	"context"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type SignedURL struct {
	FileID    string `json:"fileID"`
	SignedURL string `json:"signedURL"`
}

type SignedURLs struct {
	SignedURLs []SignedURL `json:"signedURLs"`
}

type StorageObjectRequest struct {
	FileRequests []FileRequest `json:"fileRequests"`
}

type FileRequest struct {
	ID          string `json:"id"`
	Entity      string `json:"entity"`
	Extension   string `json:"extension"`
	ContentType string `json:"contentType"`
}

type EntityBelongToFile struct {
	UserID string
}

func CreateBlankFileToGCS(ctx context.Context, fileID string, fileEntity string, fileRequest FileRequest) (string, error) {
	filePath := storageObjectFilePath(fileEntity, fileID, fileRequest.Extension)
	bucketName := infrastructure.MaterialBucketName()

	if err := infrastructure.CreateObjectToGCS(ctx, bucketName, filePath, fileRequest.ContentType, nil); err != nil {
		return "", err
	}

	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "PUT", fileRequest.ContentType)
	if err != nil {
		return "", err
	}

	return url, err
}

func EntityOfRequestedFile(ctx context.Context, entityID string, entityName string) (EntityBelongToFile, error) {
	entity := new(EntityBelongToFile)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *entity, err
	}

	key := datastore.NameKey(entityName, entityID, nil)
	if err := client.Get(ctx, key, entity); err != nil {
		return *entity, err
	}

	return *entity, nil
}

func GetSignedURL(ctx context.Context, request FileRequest) (string, error) {
	filePath := storageObjectFilePath(request.Entity, request.ID, request.Extension)
	bucketName := infrastructure.MaterialBucketName()
	url, err := infrastructure.GetGCSSignedURL(ctx, bucketName, filePath, "GET", "")
	if err != nil {
		return "", nil
	}

	return url, err
}

func storageObjectFilePath(entity string, id string, extension string) string {
	return strings.ToLower(entity) + "/" + id + "." + extension
}
