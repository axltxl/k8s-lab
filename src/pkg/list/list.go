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

// Package list implements TODO lists, it offers the Todo struct
// that is comprised of several Tasks, each one of them with
// their own Id and Message, both can be encoded as JSON strings
// through their own ToJson methods
package list

import "encoding/json"

// TodoList
////////
type TodoList struct {
	Tasks []Task `json:"tasks"`
}

func (t *TodoList) ToJson() (json_string string, err error) {
	return toJson(t)
}

// Task
////////
type Task struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (t *Task) ToJson() (json_string string, err error) {
	return toJson(t)
}

func toJson(t interface{}) (json_string string, err error) {

	json_data, err := json.Marshal(t)
	if err != nil {
		return
	}
	json_string = string(json_data)
	return
}
