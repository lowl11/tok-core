package script_event

func (event *Event) StartScript(script string) string {
	return event.startScripts[script+".sql"]
}

func (event *Event) Script(folder, script string) string {
	return event.scripts[folder].(map[string]string)[script+".sql"]
}