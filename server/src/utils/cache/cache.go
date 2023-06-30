package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_THRESHOLD = 2 << 30 // 1GB in bytes
	DEFAULT_DIRNAME   = "downloads-cache"
	CACHE_FOLDER_NAME = "cache"
	DB_FILE_NAME      = "cache_db.json"
)

var (
	setupDone      = false
	dir            string
	threshold      int64
	fileInfoMap    map[string]FileInfo // File info map in memory
	fileInfoMapMux sync.Mutex
)

// Represents information about a cached file
type FileInfo struct {
	FileName   string
	Size       int64
	Creation   int64
	LastAccess int64
}

func Setup(baseDir string, thresholdSize int64) {
	setupDone = true
	dir = baseDir
	threshold = thresholdSize

	// Create cache folder if it doesn't exist
	cacheFilesPath := filepath.Join(dir, CACHE_FOLDER_NAME)
	if err := os.MkdirAll(cacheFilesPath, 0755); err != nil {
		panic(err)
	}

	if err := initializeDb(); err != nil {
		panic(err)
	}
}

// Downloads a file from the given URL and saves it to the cache folder
func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Deletes the oldest files in the cache folder until the total size is below the threshold
func deleteOldestFiles(excludeURL string) error {
	cachePath := filepath.Join(dir, CACHE_FOLDER_NAME)

	// Create a list of file info entries
	var fileInfoList []FileInfo
	for _, fileInfo := range fileInfoMap {
		fileInfoList = append(fileInfoList, fileInfo)
	}

	// Sort the files by creation time in ascending order
	sort.Slice(fileInfoList, func(i, j int) bool {
		return fileInfoList[i].Creation < fileInfoList[j].Creation
	})

	// Calculate the total size of files in the cache folder
	totalSize := int64(0)
	for _, fileInfo := range fileInfoList {
		totalSize += fileInfo.Size
	}

	// Delete the oldest files until the total size is below the threshold,
	// excluding the just downloaded file
	for totalSize > threshold {
		oldestFile := fileInfoList[0].FileName
		oldestFilePath := filepath.Join(cachePath, oldestFile)
		oldestURL := ""

		for url, fileInfo := range fileInfoMap {
			if fileInfo.FileName == oldestFile && url != excludeURL {
				oldestURL = url
				break
			}
		}

		if oldestURL != "" {
			if err := os.Remove(oldestFilePath); err != nil {
				return err
			}
			totalSize -= fileInfoList[0].Size
			delete(fileInfoMap, oldestURL)
		}

		fileInfoList = fileInfoList[1:] // Remove the oldest file from the list
	}

	return nil
}

// Saves the file info map to the database file
func saveFileInfoMap() error {
	dbFilePath := filepath.Join(dir, DB_FILE_NAME)

	// Create a map to store the compact file info data
	fileInfoData := make(map[string][]interface{})

	// Convert file info map to compact format
	for url, fileInfo := range fileInfoMap {
		fileInfoData[url] = []interface{}{
			fileInfo.FileName,
			fileInfo.Size,
			fileInfo.Creation,
			fileInfo.LastAccess,
		}
	}

	// Serialize the compact file info data
	data, err := json.Marshal(fileInfoData)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dbFilePath, data, 0644); err != nil {
		return err
	}
	return nil
}

// Initializes the file info map from the database file
func initializeDb() error {
	dbFilePath := filepath.Join(dir, DB_FILE_NAME)

	// Check if the database file exists
	if _, err := os.Stat(dbFilePath); err == nil {
		data, err := os.ReadFile(dbFilePath)
		if err != nil {
			return err
		}

		// Deserialize the compact file info data
		var fileInfoData map[string][]interface{}
		if err := json.Unmarshal(data, &fileInfoData); err != nil {
			return err
		}

		// Convert file info data to map format
		fileInfoMap = make(map[string]FileInfo)
		for url, fileInfo := range fileInfoData {
			if len(fileInfo) != 4 {
				err = fmt.Errorf("incorrect format of file info data for URL %s", url)
				logrus.Error(err)
				return err
			}

			fileName, ok1 := fileInfo[0].(string)
			size, ok2 := fileInfo[1].(float64)
			creationUnix, ok3 := fileInfo[2].(float64)
			lastAccessUnix, ok4 := fileInfo[3].(float64)

			if !ok1 || !ok2 || !ok3 || !ok4 {
				err = fmt.Errorf("incorrect type in file info data for URL %s", url)
				logrus.Error(err)
				return err
			}

			fileInfoMap[url] = FileInfo{
				FileName:   fileName,
				Size:       int64(size),
				Creation:   int64(creationUnix),
				LastAccess: int64(lastAccessUnix),
			}
		}
	} else if os.IsNotExist(err) {
		// If the database file doesn't exist, create an empty file info map
		fileInfoMap = make(map[string]FileInfo)
	} else {
		return err
	}

	return nil
}

// Downloads a file from the given URL, stores it in the cache folder, and manages the cache size by
// deleting older files when the threshold is reached. If the file is already present in the cache,
// it will not make any HTTP requests. Returns the path to the cached file.
func DownloadAndCacheFile(url string) (string, error) {
	if !setupDone {
		defaultBaseDir := filepath.Join(os.TempDir(), DEFAULT_DIRNAME)
		Setup(defaultBaseDir, DEFAULT_THRESHOLD)
	}

	cachePath := filepath.Join(dir, CACHE_FOLDER_NAME)

	fileInfoMapMux.Lock()
	defer fileInfoMapMux.Unlock()

	// Check if the file is already cached
	if fileInfo, ok := fileInfoMap[url]; ok {
		filePath := filepath.Join(cachePath, fileInfo.FileName)
		fileInfo.LastAccess = time.Now().UnixMilli()

		// Update the last access time for the file
		fileInfoMap[url] = fileInfo
		if err := saveFileInfoMap(); err != nil {
			err = fmt.Errorf("DownloadAndCacheFile: failed to save file info map: %s", err)
			logrus.Error(err)
			return "", err
		}

		return filePath, nil
	}

	// Generate a unique file name
	fileName, err := gonanoid.Nanoid()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(cachePath, fileName)

	// Download the file
	if err := downloadFile(url, filePath); err != nil {
		return "", err
	}

	// Get file information
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	// Add the file info to the map
	fileInfoMap[url] = FileInfo{
		FileName:   fileName,
		Size:       fileInfo.Size(),
		Creation:   time.Now().UnixMilli(),
		LastAccess: time.Now().UnixMilli(),
	}

	// Check if the cache size exceeds the threshold and delete older files if necessary
	if err := deleteOldestFiles(url); err != nil {
		logrus.Errorf("Failed to delete oldest files: %v", err)
	}

	// Save the file info map to the database file
	if err := saveFileInfoMap(); err != nil {
		logrus.Errorf("Failed to save file info map: %v", err)
	}

	return filePath, nil
}

// Deletes the cached file corresponding to the given URL
func DeleteCachedFile(url string) error {
	if !setupDone {
		return fmt.Errorf("cache module is not set up")
	}

	cachePath := filepath.Join(dir, CACHE_FOLDER_NAME)

	fileInfoMapMux.Lock()
	defer fileInfoMapMux.Unlock()

	if fileInfo, ok := fileInfoMap[url]; ok {
		filePath := filepath.Join(cachePath, fileInfo.FileName)

		// Remove the file from the cache folder
		if err := os.Remove(filePath); err != nil {
			return err
		}

		// Remove the file info from the map
		delete(fileInfoMap, url)

		// Save the updated file info map to the database file
		if err := saveFileInfoMap(); err != nil {
			logrus.Errorf("Failed to save file info map: %v", err)
		}

		return nil
	}

	return fmt.Errorf("file not found in cache")
}

// Deletes all files in the cache folder and clears the file info map
func ClearCache() error {
	cachePath := filepath.Join(dir, CACHE_FOLDER_NAME)

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil
	}

	// Remove all files in the cache folder
	if err := os.RemoveAll(cachePath); err != nil {
		return err
	}

	fileInfoMapMux.Lock()
	defer fileInfoMapMux.Unlock()

	// Clear the file info map
	fileInfoMap = make(map[string]FileInfo)

	// Save the updated file info map to the database file
	if err := saveFileInfoMap(); err != nil {
		return err
	}

	return nil
}
