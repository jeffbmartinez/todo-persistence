package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
	"github.com/jeffbmartinez/todo/task"
)

/*
TaskParams is the json struct that gets passed in the request to create
a new Task object.
*/
type TaskParams struct {
	Name       string    `json:"name"`
	ParentIDs  []string  `json:"parentIDs"`
	DueDate    time.Time `json:"dueDate"`
	Categories []string  `json:"categories"`
}

// NewTask handles requests to the /tasks/new endpoint.
func NewTask(response http.ResponseWriter, request *http.Request) {
	handler := BasicResponse(http.StatusMethodNotAllowed)

	switch request.Method {
	case "POST":
		handler = postNewTask
	}

	handler(response, request)
}

func postNewTask(response http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)

	var params TaskParams
	err := decoder.Decode(&params)
	if err != nil {
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	tasklist, err := storage.GetTasklist()
	if err != nil {
		log.Errorf("Could not get tasklist")
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	var parentTasks []*task.Task
	for _, parentID := range params.ParentIDs {
		parentTask, ok := tasklist.Registry[parentID]
		if !ok {
			WriteBasicResponse(http.StatusBadRequest, response, request)
			return
		}

		parentTasks = append(parentTasks, parentTask)
	}

	newTask := tasklist.AddTask(params.Name, parentTasks)
	newTask.DueDate = params.DueDate
	newTask.Categories = params.Categories

	if err := storage.SaveTasklist(tasklist); err != nil {
		log.Errorf("Couldn't save tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteJSONResponse(response, newTask, http.StatusOK)
}
