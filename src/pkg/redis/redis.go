package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/axltxl/k8s-lab/src/pkg/config"
	"github.com/axltxl/k8s-lab/src/pkg/list"
	"github.com/go-redis/redis/v8"
)

var redisHost string = config.RedisHost
var redisPort string = config.RedisPort

var redisCtx = context.Background()

// FIXME: doc me
func connect() *redis.Client {
	log.Printf("Redis: connecting to host %s:%s", redisHost, redisPort)
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// FIXME: doc me
func PushTask(task *list.Task) {
	rdb := connect()
	rdb.Set(redisCtx, task.Id, task.Message, 0)
}

// FIXME: doc me
func GetTodoList() *list.TodoList {
	rdb := connect()

	keys, _ := rdb.Keys(redisCtx, "*").Result()
	tasks := make([]list.Task, len(keys))

	for i, task_id := range keys {
		task_message, _ := rdb.Get(redisCtx, task_id).Result()
		tasks[i] = list.Task{Id: task_id, Message: task_message}
	}

	//
	return &list.TodoList{
		Tasks: tasks,
	}
}
