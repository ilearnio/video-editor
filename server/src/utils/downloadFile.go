package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url, destPath string) error {
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
