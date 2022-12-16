package image_event

type Event struct {
	basePath string
}

func Create(basePath string) *Event {
	return &Event{
		basePath: basePath,
	}
}
