package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/axltxl/k8s-lab/src/pkg/config"
	"github.com/axltxl/k8s-lab/src/pkg/list"
	"github.com/axltxl/k8s-lab/src/pkg/redis"
	"github.com/axltxl/k8s-lab/src/pkg/uuid"
)

// FIXME: doc me
// var config.HttpPort string = config.Get().HttpPort
// var config.HttpPort string = config.HttpPort

/* todoList resource */
func todoListHandler(w http.ResponseWriter, r *http.Request) {

	// FIXME: set Content-Type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get a *list.TodoList
	todolist := redis.GetTodoList()

	// Decode it as a JSON string (actually this returns a []byte)
	json_data, _ := json.Marshal(todolist)

	// and give the answer back
	fmt.Fprintln(w, string(json_data))
}

/* todoListItem resource */
func todoListTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task *list.Task

	// request_data := make(map[string]interface{})
	// anonymous struct
	request_data := struct {
		Message string `json:"message"`
	}{
		"",
	}

	// Unmarshaling this was pain in the ass to understand
	request_body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(request_body, &request_data)

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {

		// Create new task
		task = &list.Task{
			Id:      uuid.GenerateUuid(),
			Message: request_data.Message,
		}

		// and push it
		redis.PushTask(task)

		// Decode it as a JSON string (actually this returns a []byte)
		json_data, _ := json.Marshal(task)

		// and give the answer back
		fmt.Fprintln(w, string(json_data))
	}
}

// FIXME: doc me
func setupHandlers() {
	http.HandleFunc("/todolist", todoListHandler)
	http.HandleFunc("/todolist/task", todoListTaskHandler)
}

func Start() error {

	//
	setupHandlers()

	// FIXME: doc me
	log.Printf("Starting server at port %s", config.HttpPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", config.HttpPort), nil)
}
