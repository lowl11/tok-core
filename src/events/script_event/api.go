package script_event

// StartScript получить скрипт из папки /resources/scripts/start/<файл>.sql
func (event *Event) StartScript(script string) string {
	return event.startScripts[script+".sql"]
}

// Script получить скрипт из папки /resources/script/<папка>/<файл>.sql
func (event *Event) Script(folder, script string) string {
	return event.scripts[folder].(map[string]string)[script+".sql"]
}
