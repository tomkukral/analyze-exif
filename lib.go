package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type exifData struct {
	filename string
	focal    float32
}

// read interesting exif data
func GetExif(filename string) exifData {

	out := exifData{
		filename: filename,
	}

	// open file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	ex, err := exif.Decode(f)
	if err != nil {
		log.Println(err)
		return out
	}

	focalData, err := ex.Get(exif.FocalLength)
	if err != nil {
		log.Println(err)
		return out
	}

	number, denom, _ := focalData.Rat2(0)
	focal := float32(number) / float32(denom)

	out.focal = focal

	return out
}

var photoSuffixes = [...]string{".jpg", ".JPG"}

func IsPhoto(filename string) bool {
	for _, suffix := range photoSuffixes {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}

	return false
}

func FindFiles(dir string) []string {
	out := make([]string, 0)

	err := filepath.Walk(dir, func(ph string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if IsPhoto(ph) {
			out = append(out, ph)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return out
}
