package script_event

import (
	"github.com/lowl11/lazyfile/folderapi"
)

func (event *Event) readStartScripts() error {
	files, err := folderapi.Objects("resources/scripts/start")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsFolder {
			continue
		}

		body, err := file.Read()
		if err != nil {
			return err
		}

		event.startScripts[file.Name] = string(body)
	}

	return nil
}

func (event *Event) readScripts() error {
	folders, err := folderapi.Objects("resources/scripts/")
	if err != nil {
		return err
	}

	for _, folder := range folders {
		if !folder.IsFolder {
			continue
		}

		folderMap := make(map[string]string)

		files, err := folderapi.Objects("resources/scripts/" + folder.Name)
		if err != nil {
			return err
		}

		for _, file := range files {
			body, err := file.Read()
			if err != nil {
				return err
			}

			folderMap[file.Name] = string(body)
		}

		event.scripts[folder.Name] = folderMap
	}

	return nil
}