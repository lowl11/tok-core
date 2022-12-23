package image_event

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/lowl11/lazy-collection/array"
	"image"
	"path/filepath"
)

var (
	extensions = array.NewWithList[string](
		".png", ".jpg", ".jpeg",
	)
)

// fromBase64 конвертация из base64 в байты
func (event *Event) fromBase64(buffer string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(buffer)
}

// validateImageName валидация на расширение файла
func (event *Event) validateImageName(name string) error {
	fileExtension := filepath.Ext(name)
	if extensions.Single(func(item string) bool { return item == fileExtension }) == nil {
		return errors.New("invalid extension")
	}

	return nil
}

func (event *Event) getSize(buffer []byte) (width int, height int) {
	reader := bufio.NewReader(bytes.NewReader(buffer))
	imageConfig, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0
	}

	return imageConfig.Width, imageConfig.Height
}
