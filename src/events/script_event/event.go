package script_event

type Event struct {
	startScripts map[string]string
	scripts      map[string]any
}

func Create() (*Event, error) {
	event := &Event{
		startScripts: make(map[string]string),
		scripts:      make(map[string]any),
	}

	if err := event.readStartScripts(); err != nil {
		return nil, err
	}

	if err := event.readScripts(); err != nil {
		return nil, err
	}

	return event, nil
}