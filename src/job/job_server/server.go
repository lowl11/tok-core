package job_server

import (
	"github.com/go-co-op/gocron"
	"time"
	"tok-core/src/controllers"
	"tok-core/src/controllers/feed_controller"
	"tok-core/src/controllers/post_controller"
)

type Server struct {
	scheduler *gocron.Scheduler

	post *post_controller.Controller
	feed *feed_controller.Controller
}

func Create(apiControllers *controllers.ApiControllers) *Server {
	return &Server{
		scheduler: gocron.NewScheduler(time.UTC),

		post: apiControllers.Post,
		feed: apiControllers.Feed,
	}
}
