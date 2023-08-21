package storage

import (
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
	GetFiles(docUUID string, imageUrls []ImageInfo) ([]ImageInfo, error)
	Cleanup(docUUID string) error
	GetObject(filename string, expiration int) (string, error)
}

func GetFile(storageRoot, docUUID string, imageInfo ImageInfo) (string, string, error) {
	resp, err := http.Get(imageInfo.Original.URL)
	if err != nil {
		slog.ErrorContext(common.CTX, "Error fetching image:", err.Error())
		return "", "", err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp(filepath.Join(storageRoot, docUUID, "images", "tmp"), fmt.Sprintf("*.%s", imageInfo.Original.Extension))
	if err != nil {
		slog.ErrorContext(common.CTX, "Error creating temporary file:", err.Error())
		return "", "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		slog.ErrorContext(common.CTX, "Error copying content to temporary file:", err.Error())
		return "", "", err
	}

	slog.DebugContext(common.CTX, fmt.Sprintf("%s => %s\n", imageInfo.Original.URL, tempFile.Name()))

	return tempFile.Name(), HashFile(tempFile.Name()), nil
}

func CreatePath(storageRoot, docUUID string) {
	directory := filepath.Join(storageRoot, docUUID, "images", "tmp")
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		slog.ErrorContext(common.CTX, "Error creating directory:", err)
	}
}

func HashFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		slog.ErrorContext(common.CTX, "Error opening file for hashing:", err)
		return ""
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		slog.ErrorContext(common.CTX, "Error copying content for hashing:", err)
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func StorageGetFiles(storage NotesStorage, db db.NotesDatabase, note db.Note) {
	imageUrls := GetAllImageUrls(note.Content)
	imageUrls, err := storage.GetFiles(note.ID.Hex(), imageUrls)
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
	}
	note.Content = ReplaceUrls(note.Content, imageUrls)
	db.Update(note)
}

func GetAllImageUrls(content string) []ImageInfo {
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

func ReplaceUrls(content string, imageUrls []ImageInfo) string {
	for _, item := range imageUrls {
		content = strings.ReplaceAll(content, item.Original.URL, item.Replacement)
	}
	return content
}
