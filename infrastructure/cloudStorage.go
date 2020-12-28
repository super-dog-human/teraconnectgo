package infrastructure

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

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

// CreateObjectToGCS creates object to GCS.
func CreateObjectToGCS(ctx context.Context, bucketName, filePath, contentType string, contents []byte) error {
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

// GetObjectFromGCS gets object from GCS.
func GetObjectFromGCS(ctx context.Context, bucketName, filePath string) ([]byte, error) {
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
	expire := time.Now().AddDate(1, 0, 0)

	url, err := storage.SignedURL(bucket, key, &storage.SignedURLOptions{
		GoogleAccessID: ServiceAccount(),
		SignBytes: func(b []byte) ([]byte, error) {
			resp, err := iamService.Projects.ServiceAccounts.SignBlob(
				ServiceAccount(),
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

// DeleteObjectsFromGCS deletes object in GCS.
func DeleteObjectsFromGCS(ctx context.Context, bucketName string, filePath string) error {
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

// GetPublicBackGroundImageURL returns public URL in GCS.
func GetPublicBackGroundImageURL(bucket string, id string) string {
	return "https://storage.googleapis.com/" + bucket + "/image/background/" + id + ".jpg"
}
