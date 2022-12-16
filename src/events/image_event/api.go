package image_event

import (
	"github.com/lowl11/lazyfile/fileapi"
	"github.com/lowl11/lazyfile/folderapi"
	"path/filepath"
	"tok-core/src/data/models"
)

func (event *Event) UploadAvatar(avatar *models.ImageAvatar, username string) error {
	buffer, err := event.fromBase64(avatar.Buffer)
	if err != nil {
		return err
	}

	if err = folderapi.Create(event.basePath, "profile"); err != nil {
		return err
	}

	profilePath := event.basePath + "/profile"
	if err = folderapi.Create(profilePath, username); err != nil {
		return err
	}

	usernamePath := profilePath + "/" + username
	fileName := "avatar" + filepath.Ext(avatar.Name)
	filePath := usernamePath + "/" + fileName

	// delete if exist
	if fileapi.Exist(filePath) {
		if err = fileapi.Delete(filePath); err != nil {
			return err
		}
	}

	return fileapi.Create(filePath, buffer)
}

func (event *Event) UploadWallpaper(wallpaper *models.ImageWallpaper, username string) error {
	buffer, err := event.fromBase64(wallpaper.Buffer)
	if err != nil {
		return err
	}

	if err = folderapi.Create(event.basePath, "profile"); err != nil {
		return err
	}

	profilePath := event.basePath + "/profile"
	if err = folderapi.Create(profilePath, username); err != nil {
		return err
	}

	usernamePath := profilePath + "/" + username
	fileName := "wallpaper" + filepath.Ext(wallpaper.Name)
	filePath := usernamePath + "/" + fileName

	// delete if exist
	if fileapi.Exist(filePath) {
		if err = fileapi.Delete(filePath); err != nil {
			return err
		}
	}

	return fileapi.Create(filePath, buffer)
}
