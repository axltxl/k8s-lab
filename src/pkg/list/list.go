// FIXME: doc me
package list

type Task struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type TodoList struct {
	Tasks []Task `json:"tasks"`
}
