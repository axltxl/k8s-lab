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

// Package redis is a Redis client abstraction
package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/axltxl/k8s-lab/src/pkg/config"
	"github.com/axltxl/k8s-lab/src/pkg/list"
	"github.com/go-redis/redis/v8"
)

var (
	redisHost   string        = config.RedisHost
	redisPort   string        = config.RedisPort
	redisCtx                  = context.Background()
	redisClient *redis.Client = nil
)

// Push a task on the redis cache server
func PushTask(task *list.Task) error {
	rdb := connect()
	return rdb.Set(redisCtx, task.Id, task.Message, 0).Err()
}

// Get a TodoList from redis server
func GetTodoList() (l *list.TodoList, err error) {
	rdb := connect()

	// KEYS *
	keys, err := rdb.Keys(redisCtx, "*").Result()
	if err != nil {
		return
	}

	// Generate an array of list.Task
	tasks := make([]list.Task, len(keys))
	for i, task_id := range keys {
		task_message, _ := rdb.Get(redisCtx, task_id).Result()
		tasks[i] = list.Task{Id: task_id, Message: task_message}
	}

	// Create TodoList
	l = &list.TodoList{
		Tasks: tasks,
	}

	return
}

// Connect to a redis server
func connect() *redis.Client {
	if redisClient == nil {
		log.Printf("Redis: connecting to host %s:%s", redisHost, redisPort)
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

	return redisClient
}
