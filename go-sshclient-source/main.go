package main

import (
	"flag"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	bind := flag.String("bind", "0.0.0.0:8443", "bind address")
	flag.Parse()

	e := echo.New()

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./public/template/*.template")),
	}

	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://multissh.github.io"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/", "./public")

	e.GET("/", listTermHandler)
	e.GET("/term", newTermHandler)
	e.POST("/term", createTermHandler)
	e.GET("/term/:id/data", linkTermDataHandler)
	e.POST("/term/:id/windowsize", setTermWindowSizeHandler)

	e.StartTLS(
		*bind,
		"./cert.crt",
		"./private.key",
	)
}
