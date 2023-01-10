package job_server

func (server *Server) Run() {
	server.tasks()
	server.scheduler.StartBlocking()
}

func (server *Server) RunAsync() {
	go server.Run()
}
