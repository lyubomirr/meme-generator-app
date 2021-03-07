package services

import (
	"os"
	"path"
)

var uploadsPath = os.Getenv("uploads_path")
var templateFilesPath = path.Join(uploadsPath, "templates")
var memeFilesPath = path.Join(uploadsPath, "memes")

var allowedMimeTypes = map[string]string{
	"image/png": ".png",
	"image/jpeg": ".jpg",
}

func init() {
	if _, err := os.Stat(templateFilesPath); os.IsNotExist(err) {
		err = os.MkdirAll(templateFilesPath, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(memeFilesPath); os.IsNotExist(err) {
		err = os.MkdirAll(memeFilesPath, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
}