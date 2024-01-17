package storage

import (
	"context"
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

func NewS3MinioStorage(ctx context.Context) (*S3MinioStorage, error) {
	ctx, span := common.TracerCmd.Start(ctx, "NewS3MinioStorage")

	// Initialize MinIO client
	endpoint := os.Getenv("MINIO_ADDRESS")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	region := os.Getenv("MINIO_REGION_NAME")
	bucketName := os.Getenv("S3_BUCKET")

	_, spaMinioNew := common.TracerCmd.Start(ctx, "minio.New")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Region: region,
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	spaMinioNew.End()

	appRoot, _ := os.Getwd()
	storageRoot := fmt.Sprintf("%s/storage", appRoot)

	storage := &S3MinioStorage{
		BucketName:  bucketName,
		AppRoot:     appRoot,
		StorageRoot: storageRoot,
		Client:      client,
	}

	span.End()
	return storage, nil
}

func (s *S3MinioStorage) GetFiles(ctx context.Context, docUUID string, imageUrls []ImageInfo) ([]ImageInfo, error) {
	ctx, span := common.TracerWeb.Start(ctx, "GetFiles")

	var modifiedUrls []ImageInfo

	for _, item := range imageUrls {
		ctx, spanImageUrlFile := common.TracerWeb.Start(ctx, "ImageUrlFile")
		CreatePath(ctx, s.StorageRoot, docUUID)
		localPath, fileHash, err := GetFile(ctx, s.StorageRoot, docUUID, item)
		if err != nil {
			continue
		}

		filename := fmt.Sprintf("%s/storage/images/%s.%s", docUUID, fileHash, item.Original.Extension)

		_, err = s.Client.FPutObject(ctx, s.BucketName, filename, localPath, minio.PutObjectOptions{})
		if err != nil {
			slog.InfoContext(ctx, "Error uploading object:", err)
			continue
		}

		item.Replacement = fmt.Sprintf("/storage/%s", filename)

		modifiedUrls = append(modifiedUrls, item)
		spanImageUrlFile.End()
	}

	span.End()

	return modifiedUrls, nil
}

func (s *S3MinioStorage) Cleanup(ctx context.Context, docUUID string) error {
	ctx, span := common.TracerWeb.Start(ctx, "Cleanup")

	ctx, spanListObjects := common.TracerWeb.Start(ctx, "client.ListObjects")
	objectsCh := s.Client.ListObjects(ctx, s.BucketName, minio.ListObjectsOptions{
		Prefix:    fmt.Sprintf("%s/", docUUID),
		Recursive: true,
	})
	spanListObjects.End()

	for object := range objectsCh {
		if object.Err != nil {
			slog.InfoContext(ctx, "Error listing object:", object.Err)
			continue
		}

		ctx, spanRemoveObjects := common.TracerWeb.Start(ctx, "client.RemoveObjects")
		err := s.Client.RemoveObject(ctx, s.BucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			slog.InfoContext(ctx, "Error deleting object:", object.Err)
		}
		spanRemoveObjects.End()
	}

	span.End()
	return nil
}

func (s *S3MinioStorage) GetObject(ctx context.Context, filename string, expiration int) (string, error) {
	ctx, span := common.TracerWeb.Start(ctx, "GetObject")
	presignedURL, err := s.Client.PresignedGetObject(ctx, s.BucketName, filename, time.Duration(expiration)*time.Hour, nil)
	if err != nil {
		slog.InfoContext(ctx, "Error generating presigned URL:", err)
		return "", err
	}
	span.End()
	return presignedURL.String(), nil
}
