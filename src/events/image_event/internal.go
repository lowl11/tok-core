package image_event

import (
	"encoding/base64"
	"errors"
	"github.com/lowl11/lazy-collection/array"
	"path/filepath"
)

var (
	extensions = array.NewWithList[string](
		".png", ".jpg", ".jpeg",
	)
)

func (event *Event) fromBase64(buffer string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(buffer)
}

func (event *Event) validateImageName(name string) error {
	fileExtension := filepath.Ext(name)
	if extensions.Single(func(item string) bool { return item == fileExtension }) == nil {
		return errors.New("invalid extension")
	}

	return nil
}
