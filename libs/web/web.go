package web

import (
	"embed"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/spf13/viper"
	"github.com/wasilak/notes-manager/libs/providers/storage"
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

func getPresignedURL(path string) (string, error) {
	url, err := storage.Storage.GetObject(path, 1)
	if err != nil {
		return "", err
	}

	return url, nil
}

func Init() {
	e := echo.New()
	e.Use(middleware.Recover())

	e.Use(middleware.Gzip())

	e.Use(slogecho.New(slog.Default()))
	e.Use(middleware.Recover())

	e.HideBanner = true

	t := &Template{
		templates: template.Must(template.ParseFS(getEmbededViews(views), "*.html")),
	}

	e.Renderer = t

	assetHandler := http.FileServer(getEmbededAssets(static))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

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

	e.Logger.Fatal(e.Start(viper.GetString("listen")))
}
