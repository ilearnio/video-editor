package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var VideoQuotesTableName = "videoQuotes"

// ensures that the model struct satisfies the models.Model interface
var _ models.Model = (*VideoQuote)(nil)

type VideoQuote struct {
	models.BaseModel

	Id              string `json:"id"`
	UserId          string `json:"userId" db:"userId"`
	VideoId         string `json:"videoId" db:"videoId"`
	Geometry        string `json:"geometry" db:"geometry"`
	IsHtmlEnabled   bool   `json:"isHtmlEnabled" db:"isHtmlEnabled"`
	Content         string `json:"content" db:"content"`
	Position        int    `json:"position" db:"position"`
	SelectedAudioId string `json:"selectedAudioId" db:"selectedAudioId"`

	Created types.DateTime `json:"created"`
	Updated types.DateTime `json:"updated"`
}

func (m *VideoQuote) TableName() string {
	return VideoQuotesTableName
}

func VideoQuoteNew(
	userId string,
	videoId string,
	geometry string,
	isHtmlEnabled bool,
	content string,
	position int,
	selectedAudioId string,
) *VideoQuote {
	return &VideoQuote{
		UserId:          userId,
		VideoId:         videoId,
		Geometry:        geometry,
		IsHtmlEnabled:   isHtmlEnabled,
		Content:         content,
		Position:        position,
		SelectedAudioId: selectedAudioId,
	}
}
