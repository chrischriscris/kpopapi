package scraperutils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

const BaseURL = "https://kpopping.com"
const BaseDir = "images"

// =========== String manipulation ==========

func ExtractLabel(text string) string {
	return strings.Split(text, "\n")[3][3:]
}

func getFilenameFromURLResponse(resp *http.Response) string {
	urlParts := strings.Split(resp.Request.URL.Path, "/")
	return urlParts[len(urlParts)-1]
}

// =========== Image handling ==============

func DownloadImage(url string, directory string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filename := getFilenameFromURLResponse(resp)
	fullPath := fmt.Sprintf("%s/%s.jpg", directory, filename)

	os.MkdirAll(directory, os.ModePerm)
	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func GetImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

