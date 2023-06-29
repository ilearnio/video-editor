package repositories

import (
	"fmt"
	"io"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	pocketbaseModels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/sirupsen/logrus"

	"videoeditor/src/helpers"
	"videoeditor/src/models"
)

func VideoQuoteAudioQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&models.VideoQuoteAudio{})
}

func VideoQuoteAudioFindById(
	app *pocketbase.PocketBase,
	id string,
) (*models.VideoQuoteAudio, error) {
	dao := app.Dao()

	record := &models.VideoQuoteAudio{}

	err := VideoQuoteAudioQuery(dao).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(record)

	if err != nil {
		err = fmt.Errorf("VideoQuoteAudioFindById: %v", err)
		logrus.Error(err)
		return nil, err
	}

	return record, nil
}

func VideoQuoteAudioCreate(
	app *pocketbase.PocketBase,
	videoQuoteAudio *models.VideoQuoteAudio,
	audioFileReader io.ReadCloser,
) error {
	fileName := helpers.BuildSafeFileName(videoQuoteAudio.Text, 30, "mp3")

	fileBytes, err := io.ReadAll(audioFileReader)
	if err != nil {
		err = fmt.Errorf("VideoQuoteAudioCreate: %v", err)
		logrus.Error(err)
		return err
	}

	audioFile, err := filesystem.NewFileFromBytes(fileBytes, fileName)
	if err != nil {
		err = fmt.Errorf("VideoQuoteAudioCreate: %v", err)
		logrus.Error(err)
		return err
	}

	collection, err := app.Dao().FindCollectionByNameOrId(models.VideoQuoteAudiosTableName)
	if err != nil {
		err = fmt.Errorf("VideoQuoteAudioCreate: %v", err)
		logrus.Error(err)
		return err
	}

	record := pocketbaseModels.NewRecord(collection)

	data, err := helpers.EncodeToMap(videoQuoteAudio)
	if err != nil {
		return err
	}

	form := forms.NewRecordUpsert(app, record)

	form.LoadData(data)

	form.AddFiles("audioFile", audioFile)

	if err := form.Submit(); err != nil {
		err = fmt.Errorf("VideoQuoteAudioCreate: %v", err)
		logrus.Error(err)
		return err
	}

	return nil
}

func VideoQuoteAudioGetLatestForQuote(
	app *pocketbase.PocketBase,
	videoId string,
	videoQuoteId string,
) (*models.VideoQuoteAudio, error) {
	dao := app.Dao()

	record := &models.VideoQuoteAudio{}

	err := VideoQuoteAudioQuery(dao).
		AndWhere(dbx.HashExp{"videoId": videoId, "videoQuoteId": videoQuoteId}).
		OrderBy("created DESC").
		Limit(1).
		One(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func VideoQuoteAudioFindSelectedAudios(
	app *pocketbase.PocketBase,
	videoId string,
) ([]*models.VideoQuoteAudio, error) {
	dao := app.Dao()

	records := []*models.VideoQuoteAudio{}

	err := VideoQuoteAudioQuery(dao).
		Select(models.VideoQuoteAudiosTableName+".*").
		InnerJoin(models.VideoQuotesTableName, dbx.NewExp(
			models.VideoQuotesTableName+".id = "+models.VideoQuoteAudiosTableName+".videoQuoteId"+
				" AND "+models.VideoQuotesTableName+".selectedAudioId = "+models.VideoQuoteAudiosTableName+".id",
			dbx.Params{})).
		AndWhere(dbx.HashExp{models.VideoQuoteAudiosTableName + ".videoId": videoId}).
		OrderBy(models.VideoQuotesTableName + ".position").
		All(&records)

	if err != nil {
		err = fmt.Errorf("VideoQuoteAudioFindSelectedAudios: %v", err)
		logrus.Error(err)
		return nil, err
	}

	return records, nil
}
