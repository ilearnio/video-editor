package services

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/pocketbase/pocketbase"

	"videoeditor/src/config"
	"videoeditor/src/models"
	"videoeditor/src/utils"
)

func videoQuoteAudioFilesDestinationPath(videoId string) string {
	return filepath.Join(videoFilesDestinationPath(videoId), "quotes")
}

func videoQuoteAudioFileGetUrl(audio *models.VideoQuoteAudio) string {
	return fmt.Sprint(
		config.EnvConfig.API_URL, "/files/videoQuoteAudios/", audio.Id, "/", audio.AudioFile,
	)
}

func videoQuoteAudioFileGetDownloadPath(
	audio *models.VideoQuoteAudio,
	baseDir string,
	index int,
) string {
	fileName := fmt.Sprint(index, "_", audio.AudioFile)
	return fmt.Sprint(baseDir, "/", fileName)
}

func videoQuoteAudioDownloadFile(audio *models.VideoQuoteAudio, destinationPath string) error {
	url := videoQuoteAudioFileGetUrl(audio)

	err := utils.DownloadFile(url, destinationPath, true)
	if err != nil {
		return fmt.Errorf("videoQuoteAudioDownloadFile: %v", err)
	}

	return nil
}

func videoQuoteAudioDownloadAllFiles(
	app *pocketbase.PocketBase,
	audios []*models.VideoQuoteAudio,
) (paths []string, e error) {
	if (len(audios)) == 0 {
		return nil, nil
	}

	baseDir := videoQuoteAudioFilesDestinationPath(audios[0].VideoId)

	errChan := make(chan error, len(audios))

	var wg sync.WaitGroup

	destinationPaths := make([]string, len(audios))
	for index, audio := range audios {
		destinationPath := videoQuoteAudioFileGetDownloadPath(audio, baseDir, index)
		destinationPaths[index] = destinationPath

		wg.Add(1)

		go func(audio *models.VideoQuoteAudio, destinationPath string) {
			defer wg.Done()

			err := videoQuoteAudioDownloadFile(audio, destinationPath)
			if err != nil {
				errChan <- err
			}
		}(audio, destinationPath)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return destinationPaths, nil
}
