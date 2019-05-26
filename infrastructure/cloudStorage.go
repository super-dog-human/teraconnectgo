package infrastructure

import (
	"bytes"
	"context"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
)

// CreateObjectToGCS is create object to GCS.
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

// GetObjectFromGCS is get object from GCS.
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

// GetGCSSignedURL is generate signed-URL for GCS object.
func GetGCSSignedURL(ctx context.Context, bucketName string, filePath string, method string, contentType string) (string, error) {
	account, _ := appengine.ServiceAccount(ctx)
	expire := time.Now().AddDate(1, 0, 0)

	url, signErr := storage.SignedURL(bucketName, filePath, &storage.SignedURLOptions{
		GoogleAccessID: account,
		SignBytes: func(b []byte) ([]byte, error) {
			_, signedBytes, err := appengine.SignBytes(ctx, b)
			return signedBytes, err
		},
		Method:      method,
		ContentType: contentType,
		Expires:     expire,
	})

	if signErr != nil {
		return url, signErr
	}

	return url, nil
}

// DeleteObjectsFromGCS is delete object in GCS.
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
