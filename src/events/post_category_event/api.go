package post_category_event

import "tok-core/src/data/models"

func (event *Event) Load(categories []models.PostCategoryGet) {
	event.mutex.Lock()
	defer event.mutex.Unlock()

	event.categories.Clear().PushList(categories...)
}

func (event *Event) Get(code string) *models.PostCategoryGet {
	event.mutex.Lock()
	defer event.mutex.Unlock()

	return event.categories.Single(func(item models.PostCategoryGet) bool {
		return item.Code == code
	})
}

func (event *Event) Add(code, name string) {
	event.mutex.Lock()
	defer event.mutex.Unlock()

	event.categories.Push(models.PostCategoryGet{
		Code: code,
		Name: name,
	})
}

func (event *Event) Delete(code string) {
	event.mutex.Lock()
	defer event.mutex.Unlock()

	for index, value := range event.categories.Slice() {
		if value.Code == code {
			event.categories.Remove(index)
			break
		}
	}
}
