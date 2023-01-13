package job_server

import (
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/definition"
)

func (server *Server) tasks() {
	logger := definition.Logger

	if _, err := server.scheduler.Cron("55 23 * * *").DoWithJobDetails(server.post.FillUnauthorizedJob); err != nil {
		logger.Error(err, "Fill unauthorized feed job error", layers.Job)
	}

	//if _, err := server.scheduler.Every(1).Hours().Do(server.post.FillUnauthorizedJob); err != nil {
	//	logger.Error(err, "Fill unauthorized feed job error", layers.Job)
	//}

	if _, err := server.scheduler.Cron("55 23 * * *").DoWithJobDetails(server.post.FillExploreJob); err != nil {
		logger.Error(err, "Fill explore feed job error", layers.Job)
	}

	//if _, err := server.scheduler.Every(1).Hours().Do(server.post.FillExploreJob); err != nil {
	//	logger.Error(err, "Fill explore feed job error", layers.Job)
	//}
}

func (server *Server) unauthorizedFeed() error {
	return nil
}

func (server *Server) exploreFeed() error {
	return nil
}
