package models

type ImageAvatar struct {
	Name   string `json:"Name"`
	Buffer string `json:"buffer"`
}

type ImageWallpaper struct {
	Name   string `json:"Name"`
	Buffer string `json:"buffer"`
}

type PostPicture struct {
	Name   string `json:"Name"`
	Buffer string `json:"buffer"`
}

type ImageConfig struct {
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
