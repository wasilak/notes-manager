package storage

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wasilak/notes-manager/libs/common"
)

type S3MinioStorage struct {
	BucketName  string
	AppRoot     string
	StorageRoot string
	Client      *minio.Client
}

func NewS3MinioStorage() (*S3MinioStorage, error) {
	// Initialize MinIO client
	endpoint := os.Getenv("MINIO_ADDRESS")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	region := os.Getenv("MINIO_REGION_NAME")
	bucketName := os.Getenv("S3_BUCKET")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Region: region,
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	appRoot, _ := os.Getwd()
	storageRoot := fmt.Sprintf("%s/storage", appRoot)

	storage := &S3MinioStorage{
		BucketName:  bucketName,
		AppRoot:     appRoot,
		StorageRoot: storageRoot,
		Client:      client,
	}

	return storage, nil
}

func (s *S3MinioStorage) GetFiles(docUUID string, imageUrls []ImageInfo) ([]ImageInfo, error) {
	var modifiedUrls []ImageInfo

	for _, item := range imageUrls {
		CreatePath(s.StorageRoot, docUUID)
		localPath, fileHash, err := GetFile(s.StorageRoot, docUUID, item)
		if err != nil {
			continue
		}

		filename := fmt.Sprintf("%s/storage/images/%s.%s", docUUID, fileHash, item.Original.Extension)

		_, err = s.Client.FPutObject(common.CTX, s.BucketName, filename, localPath, minio.PutObjectOptions{})
		if err != nil {
			slog.InfoContext(common.CTX, "Error uploading object:", err)
			continue
		}

		item.Replacement = fmt.Sprintf("/storage/%s", filename)

		modifiedUrls = append(modifiedUrls, item)
	}

	return modifiedUrls, nil
}

func (s *S3MinioStorage) Cleanup(docUUID string) error {
	objectsCh := s.Client.ListObjects(common.CTX, s.BucketName, minio.ListObjectsOptions{
		Prefix:    fmt.Sprintf("%s/", docUUID),
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			slog.InfoContext(common.CTX, "Error listing object:", object.Err)
			continue
		}
		err := s.Client.RemoveObject(common.CTX, s.BucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			slog.InfoContext(common.CTX, "Error deleting object:", object.Err)
		}
	}

	return nil
}

func (s *S3MinioStorage) GetObject(filename string, expiration int) (string, error) {
	presignedURL, err := s.Client.PresignedGetObject(common.CTX, s.BucketName, filename, time.Duration(expiration)*time.Hour, nil)
	if err != nil {
		slog.InfoContext(common.CTX, "Error generating presigned URL:", err)
		return "", err
	}
	return presignedURL.String(), nil
}
