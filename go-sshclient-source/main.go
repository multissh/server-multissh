package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// init
	app := echo.New()

	// handle recovers from panics anywhere
	app.Use(middleware.Recover())

	// handle CORS
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://multissh.github.io"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// static route
	app.Static("/.well-known", "./.well-known")

	// live ssh route
	app.GET("/", listTermHandler)
	app.POST("/term", createTermHandler)
	app.POST("/term/:id/windowsize", setTermWindowSizeHandler)
	app.GET("/term/:id/data", linkTermDataHandler)

	// snippet ssh route
	app.GET("/run", runCmd)

	// config route
	app.GET("/server", getConfig)
	app.GET("/snippets", getConfig)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// start http server
	go func() {
		log.Fatal(app.Start(":80"))
	}()

	// start https server
	go func() {
		log.Fatal(app.StartTLS(":443", "./cert.crt", "./private.key"))
	}()

	// signal to gracefully shutdown the server with a timeout of 10 seconds
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app.Shutdown(ctx)
}
