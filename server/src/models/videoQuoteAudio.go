package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var VideoQuoteAudiosTableName = "videoQuoteAudios"

// ensures that the model struct satisfies the models.Model interface
var _ models.Model = (*VideoQuoteAudio)(nil)

type VideoQuoteAudio struct {
	models.BaseModel

	Id           string  `json:"id"`
	UserId       string  `json:"userId" db:"userId"`
	VideoId      string  `json:"videoId" db:"videoId"`
	VideoQuoteId string  `json:"videoQuoteId" db:"videoQuoteId"`
	Text         string  `json:"text" db:"text"`
	TtsProvider  string  `json:"ttsProvider" db:"ttsProvider"`
	Voice        string  `json:"voice"`
	Seed         string  `json:"seed"`
	Speed        float64 `json:"speed"` // 0 - 5.0
	Duration     float64 `json:"duration"`
	Size         float64 `json:"size"`
	AudioFile    string  `json:"audioFile" db:"audioFile"`

	Created types.DateTime `json:"created"`
	Updated types.DateTime `json:"updated"`
}

func (m *VideoQuoteAudio) TableName() string {
	return VideoQuoteAudiosTableName
}

func VideoQuoteAudioNew(
	userId string,
	videoId string,
	videoQuoteId string,
	text string,
	ttsProvider string,
	voice string,
	seed string,
	speed float64,
	duration float64,
	size float64,
) *VideoQuoteAudio {
	return &VideoQuoteAudio{
		UserId:       userId,
		VideoId:      videoId,
		VideoQuoteId: videoQuoteId,
		Text:         text,
		TtsProvider:  ttsProvider,
		Voice:        voice,
		Seed:         seed,
		Speed:        speed,
		Duration:     duration,
		Size:         size,
	}
}
