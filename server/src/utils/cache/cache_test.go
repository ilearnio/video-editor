package cache

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils | Cache", func() {
	var testCachePath = filepath.Join(os.TempDir(), "cache-test")
	var thresholdSize int64 = 10 * 1024 * 1024 // 10Mb

	BeforeEach(func() {
		Setup(testCachePath, thresholdSize)

		// Start mocking HTTP requests
		httpmock.Activate()
	})

	AfterEach(func() {
		// Remove the temporary cache folder and test files
		Expect(os.RemoveAll(testCachePath)).To(Succeed())

		// Reset HTTP request mocking
		httpmock.DeactivateAndReset()
	})

	Context("DownloadAndCacheFile", func() {
		It("should download and cache a file", func() {
			testURL := "https://example.com/testfile.txt"

			// Mock the HTTP request for the test file
			httpmock.RegisterResponder("GET", testURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(http.StatusOK, "Test file content")
					return resp, nil
				},
			)

			// Perform the download and caching
			filePath, err := DownloadAndCacheFile(testURL)
			Expect(err).To(Succeed())
			Expect(filePath).To(BeAnExistingFile())

			// Verify that the file info is stored in the database
			fileInfoMap := readFileInfoDatabase(testCachePath)
			Expect(len(fileInfoMap)).To(Equal(1))

			fileInfo, found := fileInfoMap[testURL]
			Expect(found).To(BeTrue())
			Expect(fileInfo.Size).To(Equal(fileSize(filePath)))
		})

		It("should manage cache size by deleting older files", func() {
			// Set the threshold size to 10MB for testing
			thresholdSize := int64(10) * 1024 * 1024

			// Generate test files with larger sizes
			testFileData := generateTestFiles(5, 4*1024*1024) // Generate 5 files of 4MB each

			// Register the mock responder for the test files
			for i, data := range testFileData {
				url := "https://example.com/testfile" + strconv.Itoa(i) + ".txt"
				httpmock.RegisterResponder("GET", url,
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewBytesResponse(http.StatusOK, data)
						return resp, nil
					},
				)
			}

			// Download and cache files to exceed the threshold
			for i := range testFileData {
				url := "https://example.com/testfile" + strconv.Itoa(i) + ".txt"
				filePath, err := DownloadAndCacheFile(url)
				Expect(err).To(Succeed())
				Expect(filePath).To(BeAnExistingFile())
			}

			// Verify that the cache size is within the threshold, by also allowing the last file to be cached
			fileInfoMap := readFileInfoDatabase(testCachePath)
			var totalSize int64
			for _, info := range fileInfoMap {
				totalSize += info.Size
			}
			Expect(totalSize).To(BeNumerically("<=", thresholdSize))

			// Verify that older files are deleted
			numFilesToKeep := int(math.Floor(float64(thresholdSize) / float64(4*1024*1024)))
			fileList, err := os.ReadDir(testCachePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(fileList)).To(Equal(numFilesToKeep))
		})

		It("should delete a cached file", func() {
			testURL := "https://example.com/testfile.txt"

			// Mock the HTTP request for the test file
			httpmock.RegisterResponder("GET", testURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(http.StatusOK, "Test file content")
					return resp, nil
				},
			)

			// Perform the download and caching
			filePath, err := DownloadAndCacheFile(testURL)
			Expect(err).To(Succeed())
			Expect(filePath).To(BeAnExistingFile())

			// Delete the cached file
			err = DeleteCachedFile(testURL)
			Expect(err).To(Succeed())

			// Verify that the file is deleted from the cache
			fileInfoMap := readFileInfoDatabase(testCachePath)
			_, found := fileInfoMap[testURL]
			Expect(found).To(BeFalse())

			// Verify that the file is deleted from the disk
			_, err = os.Stat(filePath)
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("should clear the cache", func() {
			testURLs := []string{
				"https://example.com/testfile1.txt",
				"https://example.com/testfile2.txt",
				"https://example.com/testfile3.txt",
			}

			// Mock the HTTP requests for the test files
			for _, url := range testURLs {
				httpmock.RegisterResponder("GET", url,
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(http.StatusOK, "Test file content")
						return resp, nil
					},
				)
			}

			// Perform the downloads and caching
			for _, url := range testURLs {
				filePath, err := DownloadAndCacheFile(url)
				Expect(err).To(Succeed())
				Expect(filePath).To(BeAnExistingFile())
			}

			// Clear the cache
			err := ClearCache()
			Expect(err).To(Succeed())

			// Verify that the cache is empty
			fileInfoMap := readFileInfoDatabase(testCachePath)
			Expect(len(fileInfoMap)).To(Equal(0))

			// Verify that all cached files are deleted from the disk
			for _, url := range testURLs {
				filePath := getCachedFilePath(testCachePath, url)
				_, err := os.Stat(filePath)
				Expect(os.IsNotExist(err)).To(BeTrue())
			}
		})
	})
})

// generates random test files with the specified size
func generateTestFiles(count int, size int64) [][]byte {
	var fileData [][]byte
	for i := 0; i < count; i++ {
		// Generate random data for the file
		data := generateRandomData(size)
		fileData = append(fileData, data)
	}

	return fileData
}

func generateRandomData(size int64) []byte {
	data := make([]byte, size)
	_, err := rand.Read(data)
	Expect(err).NotTo(HaveOccurred())
	return data
}

// readFileInfoDatabase reads the file info database and returns the file info map
func readFileInfoDatabase(testCachePath string) map[string]FileInfo {
	dbFilePath := filepath.Join(testCachePath, DB_FILE_NAME)
	dbData, err := os.ReadFile(dbFilePath)
	Expect(err).NotTo(HaveOccurred())

	var fileInfoMap map[string][]interface{}
	Expect(json.Unmarshal(dbData, &fileInfoMap)).To(Succeed())

	convertedMap := make(map[string]FileInfo)
	for url, fileInfo := range fileInfoMap {
		fileName, ok1 := fileInfo[0].(string)
		size, ok2 := fileInfo[1].(float64)
		creationUnix, ok3 := fileInfo[2].(float64)
		lastAccessUnix, ok4 := fileInfo[3].(float64)

		if !ok1 || !ok2 || !ok3 || !ok4 {
			Fail("Invalid file info data")
		}

		convertedMap[url] = FileInfo{
			FileName:   fileName,
			Size:       int64(size),
			Creation:   int64(creationUnix),
			LastAccess: int64(lastAccessUnix),
		}
	}

	return convertedMap
}

func fileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	Expect(err).NotTo(HaveOccurred())

	return fileInfo.Size()
}

// returns the file path for a cached file based on its URL
func getCachedFilePath(testCachePath, url string) string {
	fileInfoMap := readFileInfoDatabase(testCachePath)
	fileInfo, found := fileInfoMap[url]
	if !found {
		return ""
	}
	return filepath.Join(testCachePath, fileInfo.FileName)
}
