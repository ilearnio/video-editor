package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"videoeditor/src/utils"
)

func ShotcutCalculateFileMD5(filePath string) (string, error) {
	fileStats, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	fileSize := fileStats.Size()
	chunkSize := int64(1024 * 1024) // 1 MB

	if fileSize <= chunkSize {
		return utils.CalculateFileMD5(filePath)
	} else {
		firstChunk, err := utils.CalculateFileChunkMD5(filePath, 0, chunkSize)
		if err != nil {
			return "", err
		}

		lastChunk, err := utils.CalculateFileChunkMD5(filePath, fileSize-chunkSize, chunkSize)
		if err != nil {
			return "", err
		}

		combinedHash := md5.New()
		combinedHash.Write([]byte(firstChunk + lastChunk))
		md5Hash := hex.EncodeToString(combinedHash.Sum(nil))
		return md5Hash, nil
	}
}
