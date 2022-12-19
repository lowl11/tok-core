package image_event

import (
	"github.com/google/uuid"
	"github.com/lowl11/lazyfile/fileapi"
	"github.com/lowl11/lazyfile/folderapi"
	"path/filepath"
	"strings"
	"tok-core/src/data/models"
)

/*
	UploadAvatar загрузка аватара
	Изображение принимается в base 64, далее конвертируется в байты
	Создается папка profile если такой нет
	Далее создается папка с названием юзернейма пользователя если ее нет
	Удаляются остальные файлы с именем avatar_
	Создается новый файл
*/
func (event *Event) UploadAvatar(avatar *models.ImageAvatar, username string) (string, error) {
	// валидируем расширение файла
	if err := event.validateImageName(avatar.Name); err != nil {
		return "", err
	}

	// конвертация из base64 в байты
	buffer, err := event.fromBase64(avatar.Buffer)
	if err != nil {
		return "", err
	}

	//resizedBuffer, err := image_resize.DoAvatar(avatar.Name, buffer)
	//if err != nil {
	//	return "", err
	//}

	// создается папка /profile если ее нет
	if err = folderapi.Create(event.basePath, "profile"); err != nil {
		return "", err
	}

	// создается папка /profile/<username> если ее нет
	profilePath := event.basePath + "/profile"
	if err = folderapi.Create(profilePath, username); err != nil {
		return "", err
	}

	// генерируем название и путь к файлу
	usernamePath := profilePath + "/" + username
	fileName := "avatar_" + uuid.New().String() + filepath.Ext(avatar.Name)
	filePath := usernamePath + "/" + fileName

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

	// создаем файл
	return fileName, fileapi.Create(filePath, buffer)
}

/*
	UploadWallpaper загрузка фона
	Изображение принимается в base 64, далее конвертируется в байты
	Создается папка profile если такой нет
	Далее создается папка с названием юзернейма пользователя если ее нет
	Удаляются остальные файлы с именем wallpaper_
	Создается новый файл
*/
func (event *Event) UploadWallpaper(wallpaper *models.ImageWallpaper, username string) (string, error) {
	// валидируем расширение файла
	if err := event.validateImageName(wallpaper.Name); err != nil {
		return "", err
	}

	// конвертация из base64 в байты
	buffer, err := event.fromBase64(wallpaper.Buffer)
	if err != nil {
		return "", err
	}

	//resizedBuffer, err := image_resize.DoWallpaper(wallpaper.Name, buffer)
	//if err != nil {
	//	return "", err
	//}

	// создается папка /profile если ее нет
	if err = folderapi.Create(event.basePath, "profile"); err != nil {
		return "", err
	}

	// создается папка /profile/<username> если ее нет
	profilePath := event.basePath + "/profile"
	if err = folderapi.Create(profilePath, username); err != nil {
		return "", err
	}

	// генерируем название и путь к файлу
	usernamePath := profilePath + "/" + username
	fileName := "wallpaper_" + uuid.New().String() + filepath.Ext(wallpaper.Name)
	filePath := usernamePath + "/" + fileName

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

	// создаем файл
	return fileName, fileapi.Create(filePath, buffer)
}
