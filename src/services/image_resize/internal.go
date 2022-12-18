package image_resize

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
)

func decodeImage(imageContent []byte, extension string) (image.Image, error) {
	switch extension {
	case ".png":
		src, err := png.Decode(bytes.NewReader(imageContent))
		if err != nil {
			return nil, err
		}
		return src, nil
	case ".jpg":
	case ".jpeg":
		src, err := jpeg.Decode(bytes.NewReader(imageContent))
		if err != nil {
			return nil, err
		}
		return src, nil
	}

	return nil, nil
}

func encodeImage(imageContent []byte, dst *image.RGBA, extension string) ([]byte, error) {
	var resultBuffer bytes.Buffer

	switch extension {
	case ".png":
		if err := png.Encode(&resultBuffer, dst); err != nil {
			return imageContent, err
		}
	case ".jpg":
	case ".jpeg":
		if err := jpeg.Encode(&resultBuffer, dst, nil); err != nil {
			return imageContent, err
		}
	}

	return resultBuffer.Bytes(), nil
}
