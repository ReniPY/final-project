package api

import (
	"net/http"

	"github.com/ReniPY/final-project/pkg/db"
)

const MaxTasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(MaxTasksLimit)
	if err != nil {
		writeJson(w, err)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})

}
