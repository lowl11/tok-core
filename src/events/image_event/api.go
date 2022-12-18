package image_event

import (
	"github.com/google/uuid"
	"github.com/lowl11/lazyfile/fileapi"
	"github.com/lowl11/lazyfile/folderapi"
	"path/filepath"
	"strings"
	"tok-core/src/data/models"
)

func (event *Event) UploadAvatar(avatar *models.ImageAvatar, username string) (string, error) {
	if err := event.validateImageName(avatar.Name); err != nil {
		return "", err
	}

	buffer, err := event.fromBase64(avatar.Buffer)
	if err != nil {
		return "", err
	}

	//resizedBuffer, err := image_resize.DoAvatar(avatar.Name, buffer)
	//if err != nil {
	//	return "", err
	//}

	if err = folderapi.Create(event.basePath, "profile"); err != nil {
		return "", err
	}

	profilePath := event.basePath + "/profile"
	if err = folderapi.Create(profilePath, username); err != nil {
		return "", err
	}

	usernamePath := profilePath + "/" + username
	fileName := "avatar_" + uuid.New().String() + filepath.Ext(avatar.Name)
	filePath := usernamePath + "/" + fileName

	// delete if exist
	if fileapi.Exist(filePath) {
		// удаляем остальные файлы
		objects, err := folderapi.Objects(usernamePath)
		if err != nil {
			return "", err
		}

		for _, obj := range objects {
			if obj.Name != fileName && strings.Contains(obj.Name, "avatar_") {
				if err = fileapi.Delete(obj.Path); err != nil {
					return "", err
				}
			}
		}
	}

	return fileName, fileapi.Create(filePath, buffer)
}

func (event *Event) UploadWallpaper(wallpaper *models.ImageWallpaper, username string) (string, error) {
	if err := event.validateImageName(wallpaper.Name); err != nil {
		return "", err
	}

	buffer, err := event.fromBase64(wallpaper.Buffer)
	if err != nil {
		return "", err
	}

	//resizedBuffer, err := image_resize.DoWallpaper(wallpaper.Name, buffer)
	//if err != nil {
	//	return "", err
	//}

	if err = folderapi.Create(event.basePath, "profile"); err != nil {
		return "", err
	}

	profilePath := event.basePath + "/profile"
	if err = folderapi.Create(profilePath, username); err != nil {
		return "", err
	}

	usernamePath := profilePath + "/" + username
	fileName := "wallpaper_" + uuid.New().String() + filepath.Ext(wallpaper.Name)
	filePath := usernamePath + "/" + fileName

	// delete if exist
	if fileapi.Exist(filePath) {
		// удаляем остальные файлы
		objects, err := folderapi.Objects(usernamePath)
		if err != nil {
			return "", err
		}

		for _, obj := range objects {
			if obj.Name != fileName && strings.Contains(obj.Name, "wallpaper_") {
				if err = fileapi.Delete(obj.Path); err != nil {
					return "", err
				}
			}
		}
	}

	return fileName, fileapi.Create(filePath, buffer)
}
