package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var VideosTableName = "videos"

// ensures that the model struct satisfies the models.Model interface
var _ models.Model = (*Video)(nil)

type Video struct {
	models.BaseModel

	Id                      string  `json:"id"`
	UserId                  string  `json:"userId" db:"userId"`
	Type                    string  `json:"type"`
	Status                  string  `json:"status"`
	Title                   string  `json:"title"`
	Heading                 string  `json:"heading"`
	HeadingIsHTML           bool    `json:"headingIsHTML" db:"headingIsHTML"`
	IntroImageFile          string  `json:"introImageFile" db:"introImageFile"`
	OutroImageFile          string  `json:"outroImageFile" db:"outroImageFile"`
	OutroOverlayImageFile   string  `json:"outroOverlayImageFile" db:"outroOverlayImageFile"`
	BackgroundImageFile     string  `json:"backgroundImageFile" db:"backgroundImageFile"`
	BackgroundAudioFile     string  `json:"backgroundAudioFile" db:"backgroundAudioFile"`
	BackgroundAudioVolume   string  `json:"backgroundAudioVolume" db:"backgroundAudioVolume"`
	BackgroundAudioDuration float64 `json:"backgroundAudioDuration" db:"backgroundAudioDuration"`
	GapBetweenQuotes        int     `json:"gapBetweenQuotes" db:"gapBetweenQuotes"`

	Created types.DateTime `json:"created"`
	Updated types.DateTime `json:"updated"`
}

func (m *Video) TableName() string {
	return VideosTableName
}

func VideoNew(
	userId string,
	status string,
	title string,
	heading string,
	_type string,
	introImageFile string,
	outroImageFile string,
	outroOverlayImageFile string,
	backgroundImageFile string,
	backgroundAudioFile string,
	backgroundAudioVolume string,
	backgroundAudioDuration float64,
	gapBetweenQuotes int,
) *Video {
	return &Video{
		UserId:                  userId,
		Status:                  status,
		Title:                   title,
		Heading:                 heading,
		Type:                    _type,
		IntroImageFile:          introImageFile,
		OutroImageFile:          outroImageFile,
		OutroOverlayImageFile:   outroOverlayImageFile,
		BackgroundImageFile:     backgroundImageFile,
		BackgroundAudioFile:     backgroundAudioFile,
		BackgroundAudioVolume:   backgroundAudioVolume,
		BackgroundAudioDuration: backgroundAudioDuration,
		GapBetweenQuotes:        gapBetweenQuotes,
	}
}
