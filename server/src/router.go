package src

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"videoeditor/src/controllers/videoController"
	"videoeditor/src/controllers/videoQuoteAudioController"
)

func registerRoutes(app *pocketbase.PocketBase, router *echo.Echo) {
	router.AddRoute(echo.Route{
		Method: http.MethodPost,
		Path:   "/api/videoQuoteAudio/generateFromText",
		Handler: func(c echo.Context) error {
			return videoQuoteAudioController.TextToSpeech(c, app)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
			apis.RequireRecordAuth(),
		},
	})

	router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/videos/:id/exportShotcutProjectMlt",
		Handler: func(c echo.Context) error {
			videoId := c.PathParam("id")
			return videoController.ExportShotcutProjectMlt(c, app, videoId)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
		},
	})

	router.AddRoute(echo.Route{
		Method: http.MethodGet,
		Path:   "/api/videos/:id/exportShotcutProjectFull",
		Handler: func(c echo.Context) error {
			videoId := c.PathParam("id")
			return videoController.ExportShotcutProjectFull(c, app, videoId)
		},
		Middlewares: []echo.MiddlewareFunc{
			apis.ActivityLogger(app),
		},
	})
}

func SetupRouter(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		registerRoutes(app, e.Router)
		return nil
	})
}
