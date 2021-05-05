package infrastructure

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

type SignedURL struct {
	FileID    string `json:"fileID"`
	SignedURL string `json:"signedURL"`
}

type SignedURLs struct {
	SignedURLs []SignedURL `json:"signedURLs"`
}

type StorageObjectRequest struct {
	LessonID     int64         `json:"lessonID"`
	FileRequests []FileRequest `json:"fileRequests"`
}

type FileRequest struct {
	ID          string `json:"id"`
	Entity      string `json:"entity"`
	Extension   string `json:"extension"`
	ContentType string `json:"contentType"`
}

type EntityBelongToFile struct {
	UserID int64
}

var iamService *iam.Service

func init() {
	cred, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		log.Printf("failed to initialize the Google client.\n")
		log.Printf("%v\n", errors.WithStack(err.(error)).Error())
		return
	}

	iamService, err = iam.New(cred)
	if err != nil {
		log.Printf("failed to initialize the IAM.\n")
		log.Printf("%v\n", errors.WithStack(err.(error)).Error())
		return
	}
}

// CreateFileToGCS creates object to GCS.
func CreateFileToGCS(ctx context.Context, bucketName, filePath, contentType string, contents []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	w := client.Bucket(bucketName).Object(filePath).NewWriter(ctx)
	w.ContentType = contentType
	defer w.Close()

	if len(contents) > 0 {
		if _, err := w.Write(contents); err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func CreateBlankFileToGCS(ctx context.Context, fileID string, fileEntity string, fileRequest FileRequest) (string, error) {
	filePath := StorageObjectFilePath(fileEntity, fileID, fileRequest.Extension)
	bucketName := MaterialBucketName()

	if err := CreateFileToGCS(ctx, bucketName, filePath, fileRequest.ContentType, nil); err != nil {
		return "", err
	}

	url, err := GetGCSSignedURL(ctx, bucketName, filePath, "PUT", fileRequest.ContentType)
	if err != nil {
		return "", err
	}

	return url, err
}

// GetFileFromGCS gets object from GCS.
func GetFileFromGCS(ctx context.Context, bucketName, filePath string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	r, err := client.Bucket(bucketName).Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(r); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// GetGCSSignedURL generates signed-URL for GCS object.
func GetGCSSignedURL(ctx context.Context, bucket string, key string, method string, contentType string) (string, error) {
	expire := time.Now().AddDate(0, 0, 1) // expire after a day.
	url, err := storage.SignedURL(bucket, key, &storage.SignedURLOptions{
		GoogleAccessID: ServiceAccountName(),
		SignBytes: func(b []byte) ([]byte, error) {
			resp, err := iamService.Projects.ServiceAccounts.SignBlob(
				ServiceAccountID(),
				&iam.SignBlobRequest{BytesToSign: base64.StdEncoding.EncodeToString(b)},
			).Context(ctx).Do()
			if err != nil {
				return nil, err
			}
			return base64.StdEncoding.DecodeString(resp.Signature)
		},
		Method:      method,
		ContentType: contentType,
		Expires:     expire,
	})

	if err != nil {
		return url, err
	}

	return url, nil
}

// DeleteObjectFromGCS deletes object in GCS.
func DeleteObjectFromGCS(ctx context.Context, bucketName string, filePath string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	o := client.Bucket(bucketName).Object(filePath)
	if err := o.Delete(ctx); err != nil {
		return err
	}

	return nil
}

// GetPublicBackGroundImageURL returns public image file URL in GCS.
func GetPublicBackGroundImageURL(id string) string {
	return "https://storage.googleapis.com/" + PublicBucketName() + "/image/background/" + id + ".jpg"
}

// GetPublicBackGroundMusicURL returns public audio file URL in GCS.
func GetPublicBackGroundMusicURL(id string) string {
	return "https://storage.googleapis.com/" + PublicBucketName() + "/audio/bgm/" + id + ".mp3"
}

func StorageObjectFilePath(entity string, id string, extension string) string {
	return fmt.Sprintf("%s/%s.%s", strings.ToLower(entity), id, extension)
}
