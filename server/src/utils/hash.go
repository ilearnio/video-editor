package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func CalculateMD5(data io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", err
	}

	md5Hash := hex.EncodeToString(hash.Sum(nil))
	return md5Hash, nil
}

func CalculateChunkMD5(data io.ReaderAt, start, length int64) (string, error) {
	sectionReader := io.NewSectionReader(data, start, length)
	return CalculateMD5(sectionReader)
}

func CalculateBytesMD5(data []byte) (string, error) {
	return CalculateMD5(bytes.NewReader(data))
}

func CalculateByteChunkMD5(data []byte, start, length int64) (string, error) {
	return CalculateChunkMD5(bytes.NewReader(data), start, length)
}

func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return CalculateMD5(file)
}

func CalculateFileChunkMD5(filePath string, start, length int64) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Seek(start, 0)
	if err != nil {
		return "", err
	}

	return CalculateChunkMD5(file, 0, length)
}
