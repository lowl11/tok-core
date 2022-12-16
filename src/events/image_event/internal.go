package image_event

import "encoding/base64"

func (event *Event) fromBase64(buffer string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(buffer)
}
