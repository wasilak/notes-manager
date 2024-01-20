package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/spf13/viper"
	"github.com/wasilak/notes-manager/libs/common"
	"github.com/wasilak/notes-manager/libs/providers/storage"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/metric"
)

//go:embed views/*
var views embed.FS

//go:embed static/**/*
var static embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getEmbededViews(views embed.FS) fs.FS {
	fsys, err := fs.Sub(views, "views")
	if err != nil {
		panic(err)
	}

	return fsys
}

func getEmbededAssets(static embed.FS) http.FileSystem {
	fsys, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func getPresignedURL(ctx context.Context, path string) (string, error) {
	url, err := storage.Storage.GetObject(ctx, path, 1)
	if err != nil {
		return "", err
	}

	return url, nil
}

func Init(ctx context.Context) {
	ctx, span := common.TracerCmd.Start(ctx, "WebInit")

	e := echo.New()

	if viper.GetBool("otelEnabled") {
		e.Use(otelecho.Middleware(os.Getenv("OTEL_SERVICE_NAME"), otelecho.WithSkipper(func(c echo.Context) bool {
			return strings.Contains(c.Path(), "static")
		})))
	}

	e.Use(middleware.Recover())

	e.Use(middleware.Gzip())

	e.Use(slogecho.New(slog.Default()))

	e.HideBanner = true

	ctx, spanTemplates := common.TracerCmd.Start(ctx, "Templates")
	t := &Template{
		templates: template.Must(template.ParseFS(getEmbededViews(views), "*.html")),
	}
	spanTemplates.End()

	e.Renderer = t

	ctx, spanAssets := common.TracerCmd.Start(ctx, "Assets")
	assetHandler := http.FileServer(getEmbededAssets(static))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	spanAssets.End()

	ctx, spanPaths := common.TracerCmd.Start(ctx, "Paths")
	e.GET("/storage/:path", storageEndpoint)

	e.GET("/api/list/", apiList)
	e.GET("/api/note/:uuid", apiNote)
	e.POST("/api/note/:uuid", apiNoteUpdate)
	e.DELETE("/api/note/:uuid", apiNoteDelete)
	e.PUT("/api/note/", apiNoteNew)

	e.POST("/api/ai/rewrite/", apiAIRewrite)

	e.GET("/api/tags/", apiTags)
	e.GET("/health", health)
	e.GET("/:path", index)
	e.GET("/", index)
	spanPaths.End()

	// Create an instance on a meter for the given instrumentation scope
	meter := common.MeterProvider.Meter(
		"github.com/wasilak/notes-manager",
		metric.WithInstrumentationVersion(common.GetVersion()),
	)

	var err error
	RequestCount, err = meter.Int64Counter(
		fmt.Sprintf("%s_request_count", os.Getenv("OTEL_SERVICE_NAME")),
		metric.WithDescription("Incoming request count"),
		metric.WithUnit("request"),
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	span.End()
	e.Logger.Fatal(e.Start(viper.GetString("listen")))
}
