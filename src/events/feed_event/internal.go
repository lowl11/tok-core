package feed_event

const (
	explorePrefix = "explore_"
)

var (
	exploreFields = []string{
		"text", "keys",
	}
)

func (event *Event) notMyAccount(username string) map[string]any {
	return map[string]any{
		"term": map[string]any{
			"author": username,
		},
	}
}
