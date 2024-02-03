package storage

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"log/slog"

	"github.com/wasilak/notes-manager/libs/common"
	"github.com/wasilak/notes-manager/libs/providers/db"
)

var Storage NotesStorage

type OrignalImage struct {
	URL       string
	Extension string
}

type ImageInfo struct {
	Original    OrignalImage
	Replacement string
}

type NotesStorage interface {
	GetFiles(ctx context.Context, docUUID string, imageUrls []ImageInfo) ([]ImageInfo, error)
	Cleanup(ctx context.Context, docUUID string) error
	GetObject(ctx context.Context, filename string, expiration int) (string, error)
}

func GetFile(ctx context.Context, storageRoot, docUUID string, imageInfo ImageInfo) (string, string, error) {
	ctx, span := common.TracerWeb.Start(ctx, "GetFile")

	ctx, spanHttpGet := common.TracerWeb.Start(ctx, "http.Get")
	resp, err := http.Get(imageInfo.Original.URL)
	if err != nil {
		common.HandleError(ctx, err)
		return "", "", err
	}
	defer resp.Body.Close()
	spanHttpGet.End()

	ctx, spanCreateTemp := common.TracerWeb.Start(ctx, "CreateTemp")
	tempFile, err := os.CreateTemp(filepath.Join(storageRoot, docUUID, "images", "tmp"), fmt.Sprintf("*.%s", imageInfo.Original.Extension))
	if err != nil {
		common.HandleError(ctx, err)
		return "", "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		common.HandleError(ctx, err)
		return "", "", err
	}
	spanCreateTemp.End()

	slog.DebugContext(ctx, fmt.Sprintf("%s => %s\n", imageInfo.Original.URL, tempFile.Name()))

	hashFile, err := HashFile(ctx, tempFile.Name())
	if err != nil {
		common.HandleError(ctx, err)
		return "", "", err
	}

	span.End()

	return tempFile.Name(), hashFile, nil
}

func CreatePath(ctx context.Context, storageRoot, docUUID string) error {
	ctx, span := common.TracerWeb.Start(ctx, "CreatePath")
	defer span.End()

	directory := filepath.Join(storageRoot, docUUID, "images", "tmp")
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		common.HandleError(ctx, err)
		return err
	}

	return nil
}

func HashFile(ctx context.Context, filename string) (string, error) {
	ctx, span := common.TracerWeb.Start(ctx, "HashFile")
	defer span.End()

	file, err := os.Open(filename)
	if err != nil {
		common.HandleError(ctx, err)
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		common.HandleError(ctx, err)
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func StorageGetFiles(ctx context.Context, storage NotesStorage, db db.NotesDatabase, note db.Note) error {
	ctx, span := common.TracerWeb.Start(ctx, "StorageGetFiles")
	defer span.End()

	imageUrls := GetAllImageUrls(ctx, note.Content)
	imageUrls, err := storage.GetFiles(ctx, note.ID.Hex(), imageUrls)
	if err != nil {
		common.HandleError(ctx, err)
		return err
	}
	note.Content = ReplaceUrls(ctx, note.Content, imageUrls)

	err = db.Update(ctx, note)
	if err != nil {
		common.HandleError(ctx, err)
		return err
	}

	return nil
}

func GetAllImageUrls(ctx context.Context, content string) []ImageInfo {
	_, span := common.TracerWeb.Start(ctx, "GetAllImageUrls")
	defer span.End()

	pattern := regexp.MustCompile(`(https?:[\/\.\w\s\-\*]*\.(jpg|gif|png|jpeg|webp|svg))`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	uniqueMatches := make(map[string]bool)
	var result []ImageInfo

	for _, match := range matches {

		orignalImage := OrignalImage{
			URL:       match[1],
			Extension: match[2],
		}

		imageInfo := ImageInfo{
			Original:    orignalImage,
			Replacement: "",
		}

		// Avoid adding duplicate URLs
		if _, exists := uniqueMatches[imageInfo.Original.URL]; !exists {
			uniqueMatches[imageInfo.Original.URL] = true
			result = append(result, imageInfo)
		}
	}

	return result
}

func ReplaceUrls(ctx context.Context, content string, imageUrls []ImageInfo) string {
	_, span := common.TracerWeb.Start(ctx, "ReplaceUrls")
	defer span.End()

	for _, item := range imageUrls {
		content = strings.ReplaceAll(content, item.Original.URL, item.Replacement)
	}
	return content
}
