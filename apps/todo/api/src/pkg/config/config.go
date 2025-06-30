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

// Package config sets up environment variable
// based configuration by reading a .env file
// at the working directory (if exists).
// Configuration variables are the following:
// - TODO_HTTP_PORT  - listening port
// - TODO_REDIS_HOST - redis cache location
// - TODO_REDIS_PORT - redis cache port
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const CFG_DEFAULT_HTTP_PORT = "8000"
const CFG_DEFAULT_REDIS_PORT = "6379"
const CFG_DEFAULT_REDIS_HOST = "localhost"

var (
	// HttpPort is the server listening port
	HttpPort string

	// RedisHost is the Redis server address
	RedisHost string

	// RedisPort is the Redis server port
	RedisPort string
)

// Get environment variable. Returns fallback
// if it's not present
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// Load .env on current working directory
func loadDotEnvFile() {
	cwd, _ := os.Getwd()
	if err := godotenv.Load(fmt.Sprintf("%s/.env", cwd)); err != nil {
		log.Printf("Couldn't load .env file at %s", cwd)
	} else {
		log.Printf("Successfuly loaded dotenv file at: %s", cwd)
	}
}

// Log current configuration state
func logConfig() {
	log.Println("Current configuration")
	log.Println("---------------------")
	log.Printf("TODO_HTTP_PORT  = %s", HttpPort)
	log.Printf("TODO_REDIS_HOST = %s", RedisHost)
	log.Printf("TODO_REDIS_PORT = %s", RedisPort)
	log.Println("---------------------")
}

func init() {
	loadDotEnvFile()

	// Set up configration variables
	HttpPort = getEnv("TODO_HTTP_PORT", CFG_DEFAULT_HTTP_PORT)
	RedisHost = getEnv("TODO_REDIS_HOST", CFG_DEFAULT_REDIS_HOST)
	RedisPort = getEnv("TODO_REDIS_PORT", CFG_DEFAULT_REDIS_PORT)

	// Log current configuration
	logConfig()
}
