package modules

import (
	"encoding/json"
	"log"
	"os"
)

const configFilename = "conf.json"

// Config for app
type Config struct {
	Key string
}

// GetAPI return api key for VT
// Return api key from config file
func GetApiKey() string {
	file, err := os.Open(configFilename)
	if err != nil {
		log.Fatal("error open config file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	errDecoder := decoder.Decode(&config)
	if errDecoder != nil {
		log.Fatal("error decode config ")
	}
	return config.Key
}
