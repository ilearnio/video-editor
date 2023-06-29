package repositories

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/sirupsen/logrus"

	"videoeditor/src/models"
)

func VideoQuoteQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&models.VideoQuote{})
}

func VideoQuoteFindById(app *pocketbase.PocketBase, id string) (*models.VideoQuote, error) {
	dao := app.Dao()

	record := &models.VideoQuote{}

	err := VideoQuoteQuery(dao).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(record)

	if err != nil {
		err = fmt.Errorf("VideoQuoteFindById: %v", err)
		logrus.Error(err)
		return nil, err
	}

	return record, nil
}

func VideoQuotesFindAllForVideo(
	app *pocketbase.PocketBase,
	videoId string,
	removeNewItemPlaceholder bool,
) ([]*models.VideoQuote, error) {
	dao := app.Dao()

	records := []*models.VideoQuote{}

	err := VideoQuoteQuery(dao).
		AndWhere(dbx.HashExp{"videoId": videoId}).
		OrderBy("position").
		All(&records)

	if err != nil {
		err = fmt.Errorf("VideoQuotesFindAllForVideo: %v", err)
		logrus.Error(err)
		return nil, err
	}

	// remove last record if it's content is empty, as it's a placeholder for a new quote
	if removeNewItemPlaceholder && len(records) > 0 && records[len(records)-1].Content == "" {
		records = records[:len(records)-1]
	}

	return records, nil
}
