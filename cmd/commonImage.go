package cmd

// A common module to get a screen image of the current machine

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/kbinani/screenshot"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
)

// Just for testing:
// gets a Base64 encoded jpeg form a file
func loadImage(filename string) string {
	fimg, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer fimg.Close()

	reader := bufio.NewReader(fimg)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	result := fmt.Sprintf("data:image/jpeg;base64, %s", encoded)
	return result
}

// takes a screenshot and returns the image
func takeShot() (*image.RGBA, error) {
	n := screenshot.NumActiveDisplays()
	if n < 0 {
		return nil, errors.New("no screen detected")
	}
	bounds := image.Rectangle{image.Point{0, 0}, image.Point{600, 800}} //
	//screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	return img, err
}

// just for testing:
// stores an image to a jpeg on disk
func storeImage(img *image.RGBA, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()
	jpeg.Encode(file, img, nil)

	fmt.Printf("File Created \"%s\"\n", fileName)
}
