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
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
)

func health(c echo.Context) error {
	_, span := common.TracerWeb.Start(c.Request().Context(), "RouteHealth")

	// Record measurements
	attrsServer := httpconv.ServerRequest("", c.Request())
	attrsClient := httpconv.ClientRequest(c.Request())
	RequestCount.Add(c.Request().Context(), 1, metric.WithAttributes(attrsServer...), metric.WithAttributes(attrsClient...))

	span.End()
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}

func index(c echo.Context) error {
	_, span := common.TracerWeb.Start(c.Request().Context(), "RouteIndex")
	span.End()
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{"app_version": common.GetVersion()})
}

func storageEndpoint(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteStorageEndpoint")
	presignedURL, err := getPresignedURL(ctx, c.Param("path"))
	if err != nil {
		span.End()
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "Error", "error": err})
	}
	span.End()
	return c.Redirect(http.StatusSeeOther, presignedURL)
}

func apiList(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiList")

	filter := c.QueryParam("filter")
	sort := c.QueryParam("sort")
	tags := c.QueryParam("tags")

	notes, err := db.DB.List(ctx, strings.ToLower(filter), strings.ToLower(sort), strings.Split(tags, ","))
	if err != nil {
		span.End()
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "Error", "error": err})
	}

	span.End()
	return c.JSON(http.StatusOK, notes)
}

func apiNote(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiNote")
	uuid := c.Param("uuid")

	note, err := db.DB.Get(ctx, uuid)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		span.End()
		return c.JSON(http.StatusNotFound, note)
	}

	span.End()
	return c.JSON(http.StatusOK, note)
}

func apiNoteUpdate(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiNoteUpdate")
	var note db.Note

	if err := c.Bind(&note); err != nil {
		span.End()
		return err
	}

	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()

	note.Updated = int(unixTimestamp)

	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.StorageGetFiles(ctx, storage.Storage, db.DB, note)
	}

	db.DB.Update(ctx, note)

	span.End()
	return c.JSON(http.StatusOK, note)
}

func apiNoteDelete(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiNoteDelete")
	uuid := c.Param("uuid")

	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.Storage.Cleanup(ctx, uuid)
		span.End()
	}

	note, err := db.DB.Delete(ctx, uuid)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		span.End()
		return c.JSON(http.StatusInternalServerError, err)
	}

	span.End()
	return c.JSON(http.StatusOK, note)
}

func apiNoteNew(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiNoteNew")
	var note db.Note

	if err := c.Bind(&note); err != nil {
		span.End()
		return err
	}

	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()

	note.Created = int(unixTimestamp)
	note.Updated = int(unixTimestamp)

	createdNote, err := db.DB.Create(ctx, note)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		span.End()
		return c.JSON(http.StatusInternalServerError, err)
	}

	// # note first has to be created, in  order to have it's ID/_id
	// # and afterwards images will have to be parsed and downloaded
	// # and note itself - updated.
	if viper.GetString("STORAGE_PROVIDER") != "none" {
		go storage.StorageGetFiles(ctx, storage.Storage, db.DB, createdNote)
	}

	span.End()
	return c.JSON(http.StatusOK, createdNote)
}

func apiAIRewrite(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiAIRewrite")

	var note db.Note

	if err := c.Bind(&note); err != nil {
		return err
	}

	updatedNote, err := openai.GetAIResponse(ctx, note)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		span.End()
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"original":  note,
		"rewritten": updatedNote,
	}

	span.End()
	return c.JSON(http.StatusOK, response)
}

func apiTags(c echo.Context) error {
	ctx, span := common.TracerWeb.Start(c.Request().Context(), "RouteApiTags")
	filter := c.QueryParam("query")

	tags, err := db.DB.Tags(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		span.End()
		return c.JSON(http.StatusInternalServerError, tags)
	}

	var filteredTags []string
	for _, v := range tags {
		if strings.Contains(v, filter) {
			filteredTags = append(filteredTags, v)
		}

	}

	span.End()
	return c.JSON(http.StatusOK, filteredTags)
}
