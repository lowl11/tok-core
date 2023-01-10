package job

import (
	"tok-core/src/controllers"
	"tok-core/src/job/job_server"
)

func RunAsync(apiControllers *controllers.ApiControllers) {
	client := job_server.Create(apiControllers)
	client.RunAsync()
}

func Run(apiControllers *controllers.ApiControllers) {
	client := job_server.Create(apiControllers)
	client.Run()
}
