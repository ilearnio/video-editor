package services

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/sirupsen/logrus"

	"videoeditor/src/config"
	"videoeditor/src/models"
	"videoeditor/src/repositories"
	"videoeditor/src/services/shotcutProjectBuilder"
	"videoeditor/src/utils"
)

func VideoToShotcutProjectMlt(app *pocketbase.PocketBase, videoId string) (string, error) {
	video, err := repositories.VideoFindById(app, videoId)
	if err != nil {
		err = fmt.Errorf("VideoToShotcutProjectMlt: %v", err)
		logrus.Error(err)
		return "", err
	}

	if video.Type == "quotes" {
		mlt, err := videoQuotesBuildMlt(app, video, nil)
		if err != nil {
			err = fmt.Errorf("VideoToShotcutProjectMlt: %v", err)
			logrus.Error(err)
			return "", err
		}

		return mlt, nil
	}

	return "", errors.New("VideoToShotcutProjectMlt: unknown video type")
}

func videoQuotesBuildMlt(
	app *pocketbase.PocketBase,
	video *models.Video,
	videoQuoteAudios []*models.VideoQuoteAudio,
) (string, error) {
	unfilteredQuotes, err := repositories.VideoQuotesFindAllForVideo(app, video.Id, true)
	if err != nil {
		err = fmt.Errorf("videoQuotesBuildMlt: %v", err)
		logrus.Error(err)
		return "", err
	}

	unfilteredAudios := videoQuoteAudios
	if videoQuoteAudios == nil {
		unfilteredAudios, err = repositories.VideoQuoteAudioFindSelectedAudios(app, video.Id)
		if err != nil {
			err = fmt.Errorf("videoQuotesBuildMlt: %v", err)
			logrus.Error(err)
			return "", err
		}
	}

	quotes := []*models.VideoQuote{}
	audios := []*models.VideoQuoteAudio{}
	for _, audio := range unfilteredAudios {
		for _, quote := range unfilteredQuotes {
			if quote.Id == audio.VideoQuoteId {
				quotes = append(quotes, quote)
				audios = append(audios, audio)
			}
		}
	}

	quoteAudioAssets := make([]shotcutProjectBuilder.AssetAudio, len(audios))
	for index, audio := range audios {
		quoteAudioAssets[index] = shotcutProjectBuilder.AssetAudio{
			Path:     path.Join("quotes", fmt.Sprint(index, "_", audio.AudioFile)),
			Duration: time.Duration(int64(audios[index].Duration) * int64(time.Millisecond)),
		}
	}

	mltXml, err := shotcutProjectBuilder.BuildQuotesProject(
		shotcutProjectBuilder.BuildQuotesProjectParams{
			BackgroundImagePath:   video.BackgroundImageFile,
			IntroImagePath:        video.IntroImageFile,
			OutroImagePath:        video.OutroImageFile,
			OutroOverlayImagePath: video.OutroOverlayImageFile,
			HeadingIsHTML:         video.HeadingIsHTML,
			HeadingContent:        video.Heading,
			QuoteAudios:           quoteAudioAssets,
			Quotes:                quotes,
			BgMusicAudio: shotcutProjectBuilder.AssetAudio{
				Path:     video.BackgroundAudioFile,
				Duration: time.Duration(int64(video.BackgroundAudioDuration) * int64(time.Millisecond)),
			},
		},
	)
	if err != nil {
		err = fmt.Errorf("videoQuotesBuildMlt: unable to build mlt project. %v", err)
		logrus.Error(err)
		return "", err
	}

	return mltXml, nil
}

func VideoToShotcutProject(app *pocketbase.PocketBase, videoId string) (string, error) {
	video, err := repositories.VideoFindById(app, videoId)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: %v", err)
		logrus.Error(err)
		return "", err
	}

	if video.Type == "quotes" {
		zipFilePath, err := videoQuotesToShotcutProjectFull(app, video)
		if err != nil {
			err = fmt.Errorf("VideoToShotcutProject: %v", err)
			logrus.Error(err)
			return "", err
		}

		return zipFilePath, nil
	}

	return "", errors.New("VideoToShotcutProject: unknown video type")
}

func videoFileGetUrl(video *models.Video, fileName string) string {
	return fmt.Sprint(config.EnvConfig.API_URL, "/files/videos/", video.Id, "/", fileName)
}

func videoFilesDestinationPath(videoId string) string {
	return os.TempDir() + "/video_" + videoId
}

func videoFileGetDownloadPath(baseDir, fileName string) string {
	fullPath := fmt.Sprint(baseDir, "/", fileName)
	os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	return fullPath
}

func videoDownloadFile(video *models.Video, baseDir, fileName string) error {
	url := videoFileGetUrl(video, fileName)
	destinationPath := videoFileGetDownloadPath(baseDir, fileName)

	err := utils.DownloadFile(url, destinationPath)
	if err != nil {
		return fmt.Errorf("videoDownloadFile: %v", err)
	}

	return nil
}

func videoDownloadAllFiles(app *pocketbase.PocketBase, video *models.Video) (e error) {
	fileNames := []string{
		video.IntroImageFile,
		video.OutroImageFile,
		video.OutroOverlayImageFile,
		video.BackgroundImageFile,
		video.BackgroundAudioFile,
		video.BackgroundAudioFile,
	}

	errChan := make(chan error, len(fileNames))

	var wg sync.WaitGroup

	baseDir := videoFilesDestinationPath(video.Id)
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		wg.Add(1)

		go func(fileName string) {
			defer wg.Done()

			err := videoDownloadFile(video, baseDir, fileName)
			if err != nil {
				errChan <- err
			}
		}(fileName)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func videoQuotesToShotcutProjectFull(app *pocketbase.PocketBase, video *models.Video) (string, error) {
	destPath := videoFilesDestinationPath(video.Id)

	if err := videoDownloadAllFiles(app, video); err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: %v", err)
		logrus.Error(err)
		return "", err
	}

	quotesAudioDestPath := videoQuoteAudioFilesDestinationPath(video.Id)
	if err := os.MkdirAll(quotesAudioDestPath, os.ModePerm); err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: %v", err)
		logrus.Error(err)
		return "", err
	}

	audios, err := repositories.VideoQuoteAudioFindSelectedAudios(app, video.Id)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: %v", err)
		logrus.Error(err)
		return "", err
	}

	_, err = videoQuoteAudioDownloadAllFiles(app, audios)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: unable to download files. %v", err)
		logrus.Error(err)
		return "", err
	}

	mltXml, err := videoQuotesBuildMlt(app, video, audios)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: unable to build mlt project. %v", err)
		logrus.Error(err)
		return "", err
	}

	// Write the mltXml to a file
	mltXmlPath := destPath + "/project.mlt"
	err = os.WriteFile(mltXmlPath, []byte(mltXml), 0644)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: unable to write to a file. %v", err)
		logrus.Error(err)
		return "", err
	}

	zipPath := destPath + ".zip"

	err = utils.ZipFolder(destPath, zipPath)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: unable to zip. %v", err)
		logrus.Error(err)
		return "", err
	}

	err = os.RemoveAll(destPath)
	if err != nil {
		err = fmt.Errorf("VideoQuotesToShotcutProject: unable to remove a folder. %v", err)
		logrus.Error(err)
		return "", err
	}

	return zipPath, nil
}
