package videoController

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"

	"videoeditor/src/services"
)

func ExportShotcutProjectMlt(c echo.Context, app *pocketbase.PocketBase, videoId string) error {
	mltContents, err := services.VideoToShotcutProjectMlt(app, videoId)
	if err != nil {
		return fmt.Errorf("ExportShotcutProject: %v", err)
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename=project.mlt")
	c.Response().Header().Set("Content-Type", "application/xml")

	// Copy the file contents to the response
	reader := bytes.NewReader([]byte(mltContents))
	_, err = io.Copy(c.Response().Writer, reader)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to write the file to the response")
	}

	return nil
}

func ExportShotcutProjectFull(c echo.Context, app *pocketbase.PocketBase, videoId string) error {
	zipFilePath, err := services.VideoToShotcutProject(app, videoId)
	if err != nil {
		return fmt.Errorf("ExportShotcutProject: %v", err)
	}

	zipFileDir := filepath.Dir(zipFilePath)
	zipFileName := filepath.Base(zipFilePath)

	c.Response().Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", zipFileName),
	)

	result := c.FileFS(zipFileName, os.DirFS(zipFileDir))

	if err = os.Remove(zipFilePath); err != nil {
		return fmt.Errorf("ExportShotcutProject: unable to remove the file. %v", err)
	}

	return result
}
