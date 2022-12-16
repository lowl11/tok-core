package user_controller

import "strings"

func (controller *Controller) preprocessUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
