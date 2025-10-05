package api

import (
	"go1f/pkg/server"
)

func Init() {
	if server.Router == nil {
		server.InitRouter()
	}
	r := server.Router

	r.Post("/api/signin", SigninHandler)

	r.Get("/api/nextdate", auth(nextDayHandler))
	r.Post("/api/task", auth(AddTaskHandler))
	r.Get("/api/tasks", auth(tasksHandler))

	r.Get("/api/task", auth(GetTaskHandler))
	r.Put("/api/task", auth(UpdateTaskHandler))

	r.Post("/api/task/done", auth(DoneTaskHandler))
	r.Delete("/api/task", auth(DeleteTaskHandler))
}
