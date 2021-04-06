// FIXME: doc me
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

var HttpPort string
var RedisHost string
var RedisPort string

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func loadDotEnvFile() {
	// FIXME: doc me
	cwd, _ := os.Getwd()
	if err := godotenv.Load(fmt.Sprintf("%s/.env", cwd)); err != nil {
		log.Printf("Couldn't load .env file at %s", cwd)
	} else {
		log.Printf("Successfuly loaded dotenv file at: %s", cwd)
	}
}

func logConfig() {
	log.Println("Current configuration")
	log.Println("---------------------")
	log.Printf("TODO_HTTP_PORT  = %s", HttpPort)
	log.Printf("TODO_REDIS_HOST = %s", RedisHost)
	log.Printf("TODO_REDIS_PORT = %s", RedisPort)
	log.Println("---------------------")
}

func init() {
	// FIXME: doc me
	loadDotEnvFile()

	// FIXME: doc me
	HttpPort = getEnv("TODO_HTTP_PORT", CFG_DEFAULT_HTTP_PORT)

	RedisHost = getEnv("TODO_REDIS_HOST", CFG_DEFAULT_REDIS_HOST)
	RedisPort = getEnv("TODO_REDIS_PORT", CFG_DEFAULT_REDIS_PORT)

	logConfig()
}
