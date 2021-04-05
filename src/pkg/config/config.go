// FIXME: doc me
package config

import "os"

const HTTP_PORT = "8000"

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

var HttpPort string = getEnv("TODO_HTTP_PORT", HTTP_PORT)
