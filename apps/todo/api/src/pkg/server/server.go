// Copyright (c) 2021 Alejandro Ricoveri
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package server implements the web server
// and it handles requests on the TodoList
// using a Redis server as storage
package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/axltxl/k8s-lab/apps/todo/api/src/pkg/config"
	"github.com/axltxl/k8s-lab/apps/todo/api/src/pkg/list"
	"github.com/axltxl/k8s-lab/apps/todo/api/src/pkg/redis"
	"github.com/axltxl/k8s-lab/apps/todo/api/src/pkg/uuid"
)

// Start web server
func Start() error {
	log.Printf("Starting server at port %s", config.HttpPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", config.HttpPort), nil)
}

// Function handler (decorators a la golang)
func standardHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Set Content-Type to JSON

		err := f(w, r)

		// Send a 500 if there's an error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

// GET /todolist
func todoListHandler(w http.ResponseWriter, r *http.Request) (err error) {

	todolist, err := redis.GetTodoList()
	if err != nil {
		return
	}

	todolist_json, err := todolist.ToJson()
	if err != nil {
		return
	}

	fmt.Fprintln(w, todolist_json)

	return
}

// Get a list.Task from request
func getNewTaskFromReq(r *http.Request) (task *list.Task, err error) {

	//{
	//	"message": "<string>"
	//}
	var request_data struct {
		Message string `json:"message"`
	}

	// Unmarshaling this was pain in the ass to understand
	request_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(request_body, &request_data); err != nil {
		return
	}

	// Create new task
	task = &list.Task{
		Id:      uuid.GenerateUuid(),
		Message: request_data.Message,
	}

	return
}

// POST /todolist/task
func todoListTaskHandler(w http.ResponseWriter, r *http.Request) (err error) {

	if r.Method == "POST" {

		var task *list.Task
		task, err = getNewTaskFromReq(r)
		if err != nil {
			return
		}

		// and push it
		if err = redis.PushTask(task); err != nil {
			return
		}

		var task_json string
		task_json, err = task.ToJson()
		if err != nil {
			return
		}

		// and give the answer back
		fmt.Fprintln(w, task_json)
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	return
}

func init() {
	http.HandleFunc("/todolist", standardHandler(todoListHandler))
	http.HandleFunc("/todolist/task", standardHandler(todoListTaskHandler))
}
