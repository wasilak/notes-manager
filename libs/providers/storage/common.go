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
	ctx, span := common.Tracer.Start(ctx, "GetFile")

	ctx, spanHttpGet := common.Tracer.Start(ctx, "http.Get")
	resp, err := http.Get(imageInfo.Original.URL)
	if err != nil {
		slog.ErrorContext(ctx, "Error fetching image", "error", err.Error())
		return "", "", err
	}
	defer resp.Body.Close()
	spanHttpGet.End()

	ctx, spanCreateTemp := common.Tracer.Start(ctx, "CreateTemp")
	tempFile, err := os.CreateTemp(filepath.Join(storageRoot, docUUID, "images", "tmp"), fmt.Sprintf("*.%s", imageInfo.Original.Extension))
	if err != nil {
		slog.ErrorContext(ctx, "Error creating temporary file", "error", err.Error())
		return "", "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Error copying content to temporary file", "error", err.Error())
		return "", "", err
	}
	spanCreateTemp.End()

	slog.DebugContext(ctx, fmt.Sprintf("%s => %s\n", imageInfo.Original.URL, tempFile.Name()))

	span.End()
	return tempFile.Name(), HashFile(ctx, tempFile.Name()), nil
}

func CreatePath(ctx context.Context, storageRoot, docUUID string) {
	ctx, spanCreatePath := common.Tracer.Start(ctx, "CreatePath")
	directory := filepath.Join(storageRoot, docUUID, "images", "tmp")
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating directory:", err)
	}
	spanCreatePath.End()
}

func HashFile(ctx context.Context, filename string) string {
	ctx, spanHashFile := common.Tracer.Start(ctx, "HashFile")
	file, err := os.Open(filename)
	if err != nil {
		slog.ErrorContext(ctx, "Error opening file for hashing:", err)
		return ""
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		slog.ErrorContext(ctx, "Error copying content for hashing:", err)
		return ""
	}

	spanHashFile.End()
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func StorageGetFiles(ctx context.Context, storage NotesStorage, db db.NotesDatabase, note db.Note) {
	ctx, spanStorageGetFiles := common.Tracer.Start(ctx, "StorageGetFiles")
	imageUrls := GetAllImageUrls(ctx, note.Content)
	imageUrls, err := storage.GetFiles(ctx, note.ID.Hex(), imageUrls)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
	note.Content = ReplaceUrls(ctx, note.Content, imageUrls)
	db.Update(ctx, note)
	spanStorageGetFiles.End()
}

func GetAllImageUrls(ctx context.Context, content string) []ImageInfo {
	_, span := common.Tracer.Start(ctx, "GetAllImageUrls")
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

	span.End()

	return result
}

func ReplaceUrls(ctx context.Context, content string, imageUrls []ImageInfo) string {
	_, span := common.Tracer.Start(ctx, "ReplaceUrls")
	for _, item := range imageUrls {
		content = strings.ReplaceAll(content, item.Original.URL, item.Replacement)
	}
	span.End()
	return content
}
