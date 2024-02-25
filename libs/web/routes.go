package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/wasilak/notes-manager/libs/common"
	"github.com/wasilak/notes-manager/libs/openai"
	"github.com/wasilak/notes-manager/libs/providers/db"
	"github.com/wasilak/notes-manager/libs/providers/storage"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{"app_version": common.GetVersion()})
}

func storageEndpoint(c echo.Context) error {
	presignedURL, err := getPresignedURL(c.Request().Context(), c.Param("path"))
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "Error", "error": err})
	}
	return c.Redirect(http.StatusSeeOther, presignedURL)
}

func apiList(c echo.Context) error {
	filter := c.QueryParam("filter")
	sort := c.QueryParam("sort")
	tags := c.QueryParam("tags")

	notes, err := db.DB.List(c.Request().Context(), strings.ToLower(filter), strings.ToLower(sort), strings.Split(tags, ","))
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "Error", "error": err})
	}
	return c.JSON(http.StatusOK, notes)
}

func apiNote(c echo.Context) error {
	uuid := c.Param("uuid")

	span := trace.SpanFromContext(c.Request().Context())

	note, err := db.DB.Get(c.Request().Context(), uuid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
		return c.JSON(http.StatusNotFound, note)
	}

	span.SetAttributes(attribute.String("note", fmt.Sprintf("%+v", note)))

	return c.JSON(http.StatusOK, note)
}

func apiNoteUpdate(c echo.Context) error {
	var note db.Note

	span := trace.SpanFromContext(c.Request().Context())

	if err := c.Bind(&note); err != nil {
		common.HandleError(c.Request().Context(), err)
		return err
	}

	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()

	note.Updated = int(unixTimestamp)

	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.StorageGetFiles(c.Request().Context(), storage.Storage, db.DB, note)
	}

	err := db.DB.Update(c.Request().Context(), note)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
		return err
	}

	span.SetAttributes(attribute.String("note", fmt.Sprintf("%+v", note)))

	return c.JSON(http.StatusOK, note)
}

func apiNoteDelete(c echo.Context) error {
	uuid := c.Param("uuid")

	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.Storage.Cleanup(c.Request().Context(), uuid)
	}

	note, err := db.DB.Delete(c.Request().Context(), uuid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, note)
}

func apiNoteNew(c echo.Context) error {
	var note db.Note

	if err := c.Bind(&note); err != nil {
		common.HandleError(c.Request().Context(), err)
		return err
	}

	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()

	note.Created = int(unixTimestamp)
	note.Updated = int(unixTimestamp)

	createdNote, err := db.DB.Create(c.Request().Context(), note)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// # note first has to be created, in  order to have it's ID/_id
	// # and afterwards images will have to be parsed and downloaded
	// # and note itself - updated.
	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.StorageGetFiles(c.Request().Context(), storage.Storage, db.DB, createdNote)
	}

	return c.JSON(http.StatusOK, createdNote)
}

func apiAIEnabled(c echo.Context) error {

	response := map[string]interface{}{
		"enabled": viper.GetBool("openAIEnabled"),
	}

	return c.JSON(http.StatusOK, response)
}

func apiAIRewrite(c echo.Context) error {
	var note db.Note

	span := trace.SpanFromContext(c.Request().Context())

	if err := c.Bind(&note); err != nil {
		common.HandleError(c.Request().Context(), err)
		return err
	}

	span.SetAttributes(attribute.String("note", fmt.Sprintf("%+v", note)))

	updatedNote, err := openai.GetAIResponse(c.Request().Context(), note)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
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

	tags, err := db.DB.Tags(c.Request().Context())
	if err != nil {
		slog.ErrorContext(c.Request().Context(), err.Error())
		common.HandleError(c.Request().Context(), err)
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
