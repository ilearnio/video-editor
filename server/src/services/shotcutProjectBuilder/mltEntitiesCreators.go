package shotcutProjectBuilder

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"videoeditor/src/models"
	"videoeditor/src/services/shotcutProjectBuilder/helpers"
)

func createChainMltEntityForAudioAsset(
	id string,
	audio AssetAudio,
	resourseDir string,
	definedDuration time.Duration,
	fadeInDuration time.Duration,
	fadeOutDuration time.Duration,
	volume string, // "-1.1" = - 1.1dB
) (Chain, error) {
	chain := Chain{
		DefinedDuration: definedDuration,
		Id:              id,
		Out:             helpers.ShotcutFormatDuration(audio.Duration),
		Children: []XMLConvertible{
			Property{"length", helpers.ShotcutFormatDuration(audio.Duration)},
			Property{"eof", "pause"},
			Property{"resource", filepath.Join(resourseDir, filepath.Base(audio.Path))},
			Property{"mlt_service", "avformat-novalidate"},
			Property{"seekable", "1"},
			Property{"audio_index", "0"},
			Property{"video_index", "-1"},
			Property{"mute_on_pause", "0"},
			Property{"video_delay", "0"},
			Property{"shotcut:caption", filepath.Base(audio.Path)},
			Property{"shotcut:producer", "avformat-novalidate"},
			Property{"xml", "was here"},
		},
	}

	if fadeOutDuration != 0 {
		fadeOutFilter := &Filter{
			Id:  "filter_" + id + "_fadeOut",
			Out: helpers.ShotcutFormatDuration(definedDuration),
			Children: []XMLConvertible{
				Property{"window", "75"},
				Property{"max_gain", "20dB"},
				Property{"level", helpers.ShotcutFormatDuration(definedDuration-fadeOutDuration) + "=0;" + helpers.ShotcutFormatDuration(definedDuration) + "=-60"},
				Property{"mlt_service", "volume"},
				Property{"shotcut:filter", "fadeOutVolume"},
				Property{"shotcut:animOut", helpers.ShotcutFormatDuration(fadeOutDuration)},
				Property{"disable", "0"},
			},
		}
		chain.Children = append(chain.Children, fadeOutFilter)
	}

	if fadeInDuration != 0 {
		fadeInFilter := Filter{
			Id:  "filter_" + id + "_fadeIn",
			Out: helpers.ShotcutFormatDuration(definedDuration),
			Children: []XMLConvertible{
				Property{"window", "75"},
				Property{"max_gain", "20dB"},
				Property{"level", "00:00:00.000=-60;" + helpers.ShotcutFormatDuration(fadeInDuration) + "=0"},
				Property{"mlt_service", "volume"},
				Property{"shotcut:filter", "fadeInVolume"},
				Property{"shotcut:animIn", helpers.ShotcutFormatDuration(fadeInDuration)},
			},
		}
		chain.Children = append(chain.Children, fadeInFilter)
	}

	if volume != "" {
		volumeFilter := Filter{
			Id:  "filter_" + id + "_volume",
			Out: helpers.ShotcutFormatDuration(definedDuration),
			Children: []XMLConvertible{
				Property{"window", "75"},
				Property{"max_gain", "20dB"},
				Property{"level", volume},
				Property{"mlt_service", "volume"},
			},
		}
		chain.Children = append(chain.Children, volumeFilter)
	}

	return chain, nil
}

func createEntryMltEntity(duration time.Duration, producerId string) PlaylistEntry {
	return PlaylistEntry{
		Producer: producerId,
		In:       "00:00:00.000",
		Out:      helpers.ShotcutFormatDuration(duration),
	}
}

func createBlankMltEntity(duration time.Duration) Blank {
	return Blank{
		Length: helpers.ShotcutFormatDuration(duration),
	}
}

func createProducerMltEntity(
	id string,
	quote *models.VideoQuote,
	duration time.Duration,
	geometry string,
	fadeInDuration time.Duration,
	fadeOutDuration time.Duration,
) Producer {
	durationStr := helpers.ShotcutFormatDuration(duration)

	var quoteHTML string
	if quote.IsHtmlEnabled {
		quoteHTML = quote.Content
	} else {
		quoteHTML = strings.ReplaceAll(template.HTMLEscapeString(quote.Content), "\n", "<br/>")
	}

	richTextFilter := Filter{
		Id:  fmt.Sprintf("producer_filter_%s", id),
		Out: durationStr,
		Children: []XMLConvertible{
			Property{"argument", ""},
			Property{"geometry", geometry},
			Property{"family", "Sans"},
			Property{"size", "48"},
			Property{"weight", "400"},
			Property{"style", "normal"},
			Property{"fgcolour", "0x000000ff"},
			Property{"bgcolour", "#00000000"},
			Property{"olcolour", "0x00000000"},
			Property{"pad", "0"},
			Property{"halign", "left"},
			Property{"valign", "top"},
			Property{"outline", "0"},
			Property{"pixel_ratio", "1"},
			Property{"mlt_service", "qtext"},
			Property{"shotcut:filter", "richText"},
			Property{
				"html",
				"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.0//EN\" \"http://www.w3.org/TR/REC-html40/strict.dtd\">\n" +
					"<html><head><meta name=\"qrichtext\" content=\"1\" /><meta charset=\"utf-8\" /><style type=\"text/css\">\n" +
					"p, li { white-space: pre-wrap; }\n" +
					"hr { height: 1px; border-width: 0; }\n" +
					"li.unchecked::marker { content: \"\\2610\"; }\n" +
					"li.checked::marker { content: \"\\2612\"; }\n" +
					"</style></head><body style=\" font-family:'.AppleSystemUIFont'; font-size:13pt; font-weight:400; font-style:normal;\">\n" +
					fmt.Sprintf("<p align=\"center\" style=\" margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px; line-height:117%%;\"><span style=\" font-family:'Mate SC'; font-size:41pt; font-weight:600; color:#ffffff;\">“%s”</p></body></html>", quoteHTML),
			},
			Property{"shotcut:animIn", "00:00:00.000"},
			Property{"shotcut:animOut", "00:00:00.000"},
			Property{"disable", "0"},
		},
	}
	fadeInFilter := Filter{
		Id:  fmt.Sprintf("filter_%s_fadeIn", id),
		Out: helpers.ShotcutFormatDuration(duration),
		Children: []XMLConvertible{
			Property{"start", "1"},
			Property{"level", "1"},
			Property{"mlt_service", "brightness"},
			Property{"shotcut:filter", "fadeInBrightness"},
			Property{"alpha", fmt.Sprintf(
				"00:00:00.000=0;%s=1",
				helpers.ShotcutFormatDuration(fadeInDuration),
			)},
			Property{"shotcut:animIn", helpers.ShotcutFormatDuration(fadeInDuration)},
		},
	}

	fadeOutFilter := Filter{
		Id:  fmt.Sprintf("filter_%s_fadeOut", id),
		Out: helpers.ShotcutFormatDuration(duration),
		Children: []XMLConvertible{
			Property{"start", "1"},
			Property{"level", "1"},
			Property{"mlt_service", "brightness"},
			Property{"shotcut:filter", "fadeOutBrightness"},
			Property{"alpha", fmt.Sprintf(
				"%s=1;%s=0",
				helpers.ShotcutFormatDuration(duration-fadeOutDuration),
				durationStr,
			)},
			Property{"shotcut:animOut", helpers.ShotcutFormatDuration(fadeOutDuration)},
		},
	}

	return Producer{
		Id:  id,
		In:  "00:00:00.000",
		Out: "03:59:59.960", // Shotcut hardcodes this value himself for all text producers
		Children: []XMLConvertible{
			// Shotcut hardcodes this value himself for all text producers
			Property{"length", "04:00:00.000"},
			Property{"eof", "pause"},
			Property{"resource", "#00000000"},
			Property{"aspect_ratio", "1"},
			Property{"mlt_service", "color"},
			Property{"mlt_image_format", "rgba"},
			Property{"shotcut:caption", "transparent"},
			Property{"shotcut:detail", "transparent"},
			Property{"ignore_points", "0"},
			Property{"xml", "was here"},
			Property{"seekable", "1"},
			richTextFilter,
			fadeInFilter,
			fadeOutFilter,
		},
	}
}
