package web

import (
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/wasilak/notes-manager/libs/common"
	"github.com/wasilak/notes-manager/libs/openai"
	"github.com/wasilak/notes-manager/libs/providers/db"
	"github.com/wasilak/notes-manager/libs/providers/storage"
)

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}

func index(c echo.Context) error {
	appVersion := common.GetAppVersion()
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{"app_version": appVersion})
}

func storageEndpoint(c echo.Context) error {
	presignedURL, err := getPresignedURL(c.Param("path"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "Error", "error": err})
	}
	return c.Redirect(http.StatusSeeOther, presignedURL)
}

func apiList(c echo.Context) error {

	filter := c.QueryParam("filter")
	sort := c.QueryParam("sort")
	tags := c.QueryParam("tags")

	notes, err := db.DB.List(strings.ToLower(filter), strings.ToLower(sort), strings.Split(tags, ","))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "Error", "error": err})
	}

	return c.JSON(http.StatusOK, notes)
}

func apiNote(c echo.Context) error {
	uuid := c.Param("uuid")

	note, err := db.DB.Get(uuid)
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
		return c.JSON(http.StatusNotFound, note)
	}

	return c.JSON(http.StatusOK, note)
}

func apiNoteUpdate(c echo.Context) error {
	var note db.Note

	if err := c.Bind(&note); err != nil {
		return err
	}

	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()

	note.Updated = int(unixTimestamp)

	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.StorageGetFiles(storage.Storage, db.DB, note)
	}

	db.DB.Update(note)

	return c.JSON(http.StatusOK, note)
}

func apiNoteDelete(c echo.Context) error {
	uuid := c.Param("uuid")

	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.Storage.Cleanup(uuid)
	}

	note, err := db.DB.Delete(uuid)
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, note)
}

func apiNoteNew(c echo.Context) error {
	var note db.Note

	if err := c.Bind(&note); err != nil {
		return err
	}

	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()

	note.Created = int(unixTimestamp)
	note.Updated = int(unixTimestamp)

	createdNote, err := db.DB.Create(note)
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	// # note first has to be created, in  order to have it's ID/_id
	// # and afterwards images will have to be parsed and downloaded
	// # and note itself - updated.
	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.StorageGetFiles(storage.Storage, db.DB, createdNote)
	}

	return c.JSON(http.StatusOK, createdNote)
}

func apiAIRewrite(c echo.Context) error {

	var note db.Note

	if err := c.Bind(&note); err != nil {
		return err
	}

	updatedNote, err := openai.GetAIResponse(note)
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"original":  note,
		"rewritten": updatedNote,
	}

	return c.JSON(http.StatusOK, response)
}

func apiTags(c echo.Context) error {
	filter := c.QueryParam("query")

	tags, err := db.DB.Tags()
	if err != nil {
		slog.ErrorContext(common.CTX, err.Error())
		return c.JSON(http.StatusInternalServerError, tags)
	}

	var filteredTags []string
	for _, v := range tags {
		if strings.Contains(v, filter) {
			filteredTags = append(filteredTags, v)
		}
	}

	return c.JSON(http.StatusOK, filteredTags)
}
