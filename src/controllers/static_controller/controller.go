package static_controller

import "tok-core/src/controllers/controller"

type Controller struct {
	controller.Base
}

func Create() *Controller {
	return &Controller{}
}