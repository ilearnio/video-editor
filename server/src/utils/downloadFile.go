package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"videoeditor/src/helpers"
	"videoeditor/src/utils/cache"
)

func DownloadFile(url, destPath string, useCache bool) error {
	if useCache {
		cacheFilePath, err := cache.DownloadAndCacheFile(url)
		if err != nil {
			return fmt.Errorf("DownloadFile: error downloading and caching file: %s", err.Error())
		}

		// If the file is in cache, copy it to the destination path
		if cacheFilePath != "" {
			if err := helpers.CopyFile(cacheFilePath, destPath); err != nil {
				return fmt.Errorf("DownloadFile: error copying cached file: %s", err.Error())
			}
			return nil
		}
	}

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching file: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DownloadFile: error fetching file: status code %d", resp.StatusCode)
	}

	// Create a local file to save the downloaded content
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("DownloadFile: error creating file: %s", err.Error())
	}

	// Copy the response body to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("DownloadFile: error copying file content: %s", err.Error())
	}

	return nil
}
