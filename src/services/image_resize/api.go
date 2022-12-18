package image_resize

import (
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"path/filepath"
)

func DoAvatar(name string, imageContent []byte) ([]byte, error) {
	src, err := decodeImage(imageContent, filepath.Ext(name))
	if err != nil {
		return nil, err
	}

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, 300, 300))

	// resize
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return encodeImage(imageContent, dst, filepath.Ext(name))
}

func DoWallpaper(name string, imageContent []byte) ([]byte, error) {
	src, _ := png.Decode(bytes.NewReader(imageContent))

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, 400, 700))

	// resize
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return encodeImage(imageContent, dst, filepath.Ext(name))
}
