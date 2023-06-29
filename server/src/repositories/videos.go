package repositories

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"

	"videoeditor/src/models"
)

func VideoQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&models.Video{})
}

func VideoFindById(app *pocketbase.PocketBase, id string) (*models.Video, error) {
	dao := app.Dao()

	record := &models.Video{}

	err := VideoQuery(dao).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(record)

	if err != nil {
		return nil, err
	}

	return record, nil
}
