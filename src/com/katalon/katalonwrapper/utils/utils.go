package utils

import (
	"log"
	"os"
)

func HandleErrorIfExists(err error, additionalMessage string) {
	if err != nil {
		if additionalMessage == "" {
			log.Fatal(err)
		}
		log.Fatalf(additionalMessage, err)
	}
}

func EnsureDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModeDir|os.ModePerm); err == nil {
		return nil
	} else {
		return err
	}
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDirectory(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}
