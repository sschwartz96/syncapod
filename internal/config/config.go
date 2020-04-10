package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config holds variables for our server
type Config struct {
	Version  float64 `json:"version,omitempty"`
	DbURI    string  `json:"db_uri,omitempty"`
	DbUser   string  `json:"db_user,omitempty"`
	DbPass   string  `json:"db_pass,omitempty"`
	Port     int     `json:"port"`
	CertFile string  `json:"cert_file"`
	KeyFile  string  `json:"key_file"`
}

// ReadConfig reads the config file encoded in JSON
func ReadConfig(path string) (*Config, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal into config var
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
