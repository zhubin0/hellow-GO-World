package thumb

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"image/jpeg"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/bmp"
)

// test of generating thumbnails from multiple sources
// SEE datamesh.com\common\global\const.go for supported formats
func TestGenThumbnail(t *testing.T) {
	// png
	buf, err := generatePng()
	assert.Nil(t, err)
	out := &bytes.Buffer{}
	err = GenThumbnail(&buf, out)
	assert.Nil(t, err)
	// jpg
	buf, err = generateJpg()
	assert.Nil(t, err)
	out = &bytes.Buffer{}
	err = GenThumbnail(&buf, out)
	assert.Nil(t, err)
	// bmp
	buf, err = generateBmp()
	assert.Nil(t, err)
	out = &bytes.Buffer{}
	err = GenThumbnail(&buf, out)
	assert.Nil(t, err)
}

func generatePng() (bytes.Buffer, error) {
	m := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{256, 256}})
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			m.SetNRGBA(x, y, color.NRGBA{uint8(x), uint8((x + y) / 2), uint8(y), 255})

		}
	}
	// Save to out.png
	buf := bytes.Buffer{}
	err := png.Encode(&buf, m)
	return buf, err
}

func generateJpg() (bytes.Buffer, error) {
	m := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{256, 256}})
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			m.SetNRGBA(x, y, color.NRGBA{uint8(x), uint8((x + y) / 2), uint8(y), 255})

		}
	}
	// Save to out.png
	buf := bytes.Buffer{}
	err := jpeg.Encode(&buf, m, nil)
	return buf, err
}

func generateBmp() (bytes.Buffer, error) {
	m := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{256, 256}})
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			m.SetNRGBA(x, y, color.NRGBA{uint8(x), uint8((x + y) / 2), uint8(y), 255})

		}
	}
	// Save to out.png
	buf := bytes.Buffer{}
	err := bmp.Encode(&buf, m)
	return buf, err
}
