package videoQuoteAudioController

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	pocketbaseModels "github.com/pocketbase/pocketbase/models"

	"videoeditor/src/models"
	"videoeditor/src/repositories"
	"videoeditor/src/services"
	"videoeditor/src/types/requests"
	"videoeditor/src/utils"
)

func TextToSpeech(c echo.Context, app *pocketbase.PocketBase) error {
	var req requests.TextToSpeechRequest
	if err := utils.JsonParseRequestBody(c, &req); err != nil {
		return err
	}

	authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*pocketbaseModels.Record)

	var videoId string

	err := app.DB().
		Select("videoId").
		From(models.VideoQuotesTableName).
		AndWhere(dbx.HashExp{"id": req.VideoQuoteId}).
		Row(&videoId)
	if err != nil {
		return err
	}

	res, err := services.PlayHtTextToSpeech(services.PlayHtTextToSpeechRequest{
		Voice: req.Voice,
		Text:  req.Text,
		Seed:  req.Seed,
		Speed: req.Speed,
	})
	if err != nil {
		return err
	}

	resp, err := http.Get(res.URL)
	if err != nil {
		return fmt.Errorf("error fetching file: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching file: server returned status code %d", resp.StatusCode)
	}

	err = repositories.VideoQuoteAudioCreate(app, models.VideoQuoteAudioNew(
		authRecord.Id,
		videoId,
		req.VideoQuoteId,
		req.Text,
		"playht",
		req.Voice,
		fmt.Sprint(req.Seed),
		req.Speed,
		float64(int64(res.Duration*1000)),
		res.Size,
	), resp.Body)
	if err != nil {
		return err
	}

	savedRecord, err := repositories.VideoQuoteAudioGetLatestForQuote(app, videoId, req.VideoQuoteId)
	if err != nil {
		return err
	}

	return c.JSON(200, savedRecord)
}
