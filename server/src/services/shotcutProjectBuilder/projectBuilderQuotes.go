package shotcutProjectBuilder

import (
	"fmt"
	"html/template"
	"math"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"videoeditor/src/models"
	"videoeditor/src/services/shotcutProjectBuilder/helpers"
	"videoeditor/src/services/shotcutProjectBuilder/templates"
)

const GAP_BETWEEN_QUOTES_DURATION = 2000 * time.Millisecond // in ms

func calcTotalVideoDuration(quoteAudios []AssetAudio) time.Duration {
	var quotesDuration int64 = 0
	for _, quoteAudio := range quoteAudios {
		quotesDuration += int64(quoteAudio.Duration)
	}
	return time.Duration(
		int64(INTRO_DURATION) +
			int64(quotesDuration) +
			int64(GAP_BETWEEN_QUOTES_DURATION)*int64(len(quoteAudios)-1) +
			int64(OUTRO_DURATION),
	)
}

func buildQuoteAudioChannelXMLNodes(quoteAudios []AssetAudio) ([]XMLNode, error) {
	xmlNodes := []XMLNode{}
	for i, quoteAudio := range quoteAudios {
		id := fmt.Sprint("chain", i)
		chain, err := createChainMltEntityForAudioAsset(id, quoteAudio, "quotes", 0, 0, 0, "")
		if err != nil {
			err = fmt.Errorf("buildQuoteAudioChannelXMLNodes: failed to create chain: %v", err)
			logrus.Error(err)
			return nil, err
		}
		xmlNode := chain.ToXMLNode()
		xmlNodes = append(xmlNodes, xmlNode)
	}

	return xmlNodes, nil
}

func buildQuoteAudioChannelsXML(quoteAudios []AssetAudio) (string, error) {
	xmlNodes, err := buildQuoteAudioChannelXMLNodes(quoteAudios)
	if err != nil {
		err = fmt.Errorf("buildQuoteAudioChannelsXML: %v", err)
		logrus.Error(err)
		return "", err
	}
	var xmlStrings []string
	for _, node := range xmlNodes {
		xmlStrings = append(xmlStrings, NodeToXML(node, 1))
	}

	result := strings.Join(xmlStrings, "\n")
	return result, nil
}

func buildQuoteAudioPlaylistEntriesXML(quoteAudios []AssetAudio) string {
	var entities []XMLConvertible

	for index, audio := range quoteAudios {
		producerId := fmt.Sprintf("chain%d", index)
		entryNode := createEntryMltEntity(audio.Duration, producerId)
		if index > 0 {
			entryNode.In = ""
		}

		entities = append(entities, entryNode)

		if index != len(quoteAudios)-1 {
			entities = append(entities, createBlankMltEntity(GAP_BETWEEN_QUOTES_DURATION))
		}
	}

	var xmlStrings []string
	for _, entity := range entities {
		xmlString := NodeToXML(entity.ToXMLNode(), 2)
		xmlStrings = append(xmlStrings, xmlString)
	}

	result := strings.Join(xmlStrings, "\n")
	return result
}

func createQuoteTextProducerEntity(
	id string,
	quote *models.VideoQuote,
	duration time.Duration,
) Producer {
	const geometry = "902 373 904 546 1"
	return createProducerMltEntity(
		id, quote, duration, geometry, 1*time.Second, 1*time.Second,
	)
}

func buildQuoteTextProducerEntriesXML(producers []Producer, quoteAudios []AssetAudio) string {
	var entity []XMLConvertible

	for index, producer := range producers {
		entity = append(entity, createEntryMltEntity(quoteAudios[index].Duration, producer.Id))

		if index != len(producers)-1 {
			entity = append(entity, createBlankMltEntity(GAP_BETWEEN_QUOTES_DURATION))
		}
	}

	var xmlStrings []string
	for _, entity := range entity {
		xmlString := NodeToXML(entity.ToXMLNode(), 2)
		xmlStrings = append(xmlStrings, xmlString)
	}

	result := strings.Join(xmlStrings, "\n")
	return result
}

func buildBgMusicNodes(bgMusicAudio AssetAudio, totalVideoDuration time.Duration) ([]Chain, error) {
	musicDuration := bgMusicAudio.Duration - BG_MUSIC_TRUNCATE_DURATION
	videoDuration := totalVideoDuration - OUTRO_DURATION

	var numberOfLoops int
	if videoDuration <= musicDuration {
		numberOfLoops = 1
	} else {
		durationDifference := float64(videoDuration) - float64(musicDuration)
		fadeDurationDiff := (float64(musicDuration) - float64(BG_MUSIC_FADE_DURATION))
		numberOfLoops = int(math.Ceil(durationDifference/fadeDurationDiff) + 1)
	}

	definedDuration := time.Duration(math.Min(float64(musicDuration), float64(totalVideoDuration)))
	lastAudioDefinedDuration := totalVideoDuration % musicDuration

	var chains []Chain
	for index := 0; index < numberOfLoops; index++ {
		var duration time.Duration
		if index == numberOfLoops-1 {
			duration = lastAudioDefinedDuration
		} else {
			duration = definedDuration
		}

		chain, err := createChainMltEntityForAudioAsset(
			fmt.Sprintf("channel_bgmusic_%d", index),
			bgMusicAudio,
			"",
			duration,
			BG_MUSIC_FADE_DURATION,
			BG_MUSIC_FADE_DURATION,
			"-11.1", // dB
		)
		if err != nil {
			return nil, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func buildBgMusicPlaylistEntriesXML(bgMusicChains []Chain, isSecondTrack bool) string {
	if len(bgMusicChains) == 0 {
		return ""
	}

	entries := make([]XMLConvertible, 0)

	for index, chain := range bgMusicChains {
		nodes := []XMLConvertible{
			createEntryMltEntity(chain.DefinedDuration, chain.Id),
		}

		if index != len(bgMusicChains)-1 {
			nodes = append(nodes, createBlankMltEntity(bgMusicChains[0].DefinedDuration-BG_MUSIC_FADE_DURATION*2))
		}

		entries = append(entries, nodes...)
	}

	if isSecondTrack {
		entries = append([]XMLConvertible{createBlankMltEntity(bgMusicChains[0].DefinedDuration - BG_MUSIC_FADE_DURATION)}, entries...)
	}

	xmlStrings := make([]string, len(entries))
	for i, entry := range entries {
		xmlStrings[i] = NodeToXML(entry.ToXMLNode(), 2)
	}

	return strings.Join(xmlStrings, "\n")
}

type BuildQuotesProjectParams struct {
	BackgroundImagePath   string
	IntroImagePath        string
	OutroImagePath        string
	OutroOverlayImagePath string
	HeadingIsHTML         bool
	HeadingContent        string
	QuoteAudios           []AssetAudio
	Quotes                []*models.VideoQuote
	BgMusicAudio          AssetAudio
}

func BuildQuotesProject(p BuildQuotesProjectParams) (string, error) {
	totalVideoDuration := calcTotalVideoDuration(p.QuoteAudios)

	// Heading

	headingDuration := helpers.ShotcutFormatDuration(
		totalVideoDuration - INTRO_DURATION - OUTRO_DURATION,
	)

	var headingHTML string
	if p.HeadingIsHTML {
		headingHTML = template.HTMLEscapeString(p.HeadingContent)
	} else {
		headingHTML = template.HTMLEscapeString(
			strings.ReplaceAll(template.HTMLEscapeString(p.HeadingContent), "\n", "<br/>"),
		)
	}

	// Quote audios

	quoteAudioChannelsXML, err := buildQuoteAudioChannelsXML(p.QuoteAudios)
	if err != nil {
		err = fmt.Errorf("BuildProject: %v", err)
		logrus.Error(err)
		return "", err
	}
	quoteAudioPlaylistEntriesXML := buildQuoteAudioPlaylistEntriesXML(p.QuoteAudios)

	// Quote texts

	quoteTextProducers := make([]Producer, len(p.Quotes))
	for index, quote := range p.Quotes {
		id := fmt.Sprintf("producer_quote_%d", index)
		duration := p.QuoteAudios[index].Duration
		quoteTextProducers[index] = createQuoteTextProducerEntity(id, quote, duration)
	}
	var quoteTextProducersXML string
	for _, producer := range quoteTextProducers {
		quoteTextProducersXML += NodeToXML(producer.ToXMLNode(), 1) + "\n"
	}
	quoteTextProducerEntriesXML := buildQuoteTextProducerEntriesXML(quoteTextProducers, p.QuoteAudios)

	// Background Music
	bgMusicChains, err := buildBgMusicNodes(p.BgMusicAudio, totalVideoDuration)
	if err != nil {
		err = fmt.Errorf("BuildProject: %v", err)
		logrus.Error(err)
		return "", err
	}

	bgMusicChainXMLs := make([]string, len(bgMusicChains))
	for index, chain := range bgMusicChains {
		bgMusicChainXMLs[index] = NodeToXML(chain.ToXMLNode(), 1)
	}

	var bgMusicChainsFirstTrack []Chain
	var bgMusicChainsSecondTrack []Chain
	var bgMusicChainsFirstTrackXML string
	var bgMusicChainsSecondTrackXML string

	for i, chain := range bgMusicChains {
		if i%2 == 0 {
			bgMusicChainsFirstTrack = append(bgMusicChainsFirstTrack, chain)
			bgMusicChainsFirstTrackXML += bgMusicChainXMLs[i] + "\n"
		} else {
			bgMusicChainsSecondTrack = append(bgMusicChainsSecondTrack, chain)
			bgMusicChainsSecondTrackXML += bgMusicChainXMLs[i] + "\n"
		}
	}
	bgMusicFirstTrackPlaylistEntriesXML := buildBgMusicPlaylistEntriesXML(bgMusicChainsFirstTrack, false)
	bgMusicSecondTrackPlaylistEntriesXML := buildBgMusicPlaylistEntriesXML(bgMusicChainsSecondTrack, true)

	params := templates.QuotesTemplateParamsNew(
		path.Base(p.BackgroundImagePath),
		path.Base(p.BackgroundImagePath),

		path.Base(p.IntroImagePath),
		path.Base(p.IntroImagePath),
		helpers.ShotcutFormatDuration(INTRO_DURATION),

		path.Base(p.OutroImagePath),
		path.Base(p.OutroImagePath),
		helpers.ShotcutFormatDuration(OUTRO_DURATION),
		helpers.ShotcutFormatDuration(
			totalVideoDuration-INTRO_DURATION-OUTRO_DURATION,
		),

		path.Base(p.OutroOverlayImagePath),
		path.Base(p.OutroOverlayImagePath),
		helpers.ShotcutFormatDuration(
			totalVideoDuration-OUTRO_OVERLAY_DURATION,
		),
		helpers.ShotcutFormatDuration(OUTRO_OVERLAY_DURATION),

		helpers.ShotcutFormatDuration(totalVideoDuration),

		headingDuration,
		headingHTML,

		quoteAudioChannelsXML,
		quoteAudioPlaylistEntriesXML,

		quoteTextProducersXML,
		quoteTextProducerEntriesXML,

		bgMusicChainsFirstTrackXML,
		bgMusicChainsSecondTrackXML,
		bgMusicFirstTrackPlaylistEntriesXML,
		bgMusicSecondTrackPlaylistEntriesXML,
	)
	return templates.CompileTemplate(templates.TemplateQuotes, params)
}
