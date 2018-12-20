package config

import (
	"encoding/json"
	"errors"
	"os"
)

type AuthConfig struct {
	PublicKey  string
	PrivateKey string
}

type Configuration struct {
	Port int
	Auth AuthConfig
}

var config *Configuration = nil

func GetConfig() (*Configuration, error) {

	if config != nil {
		return config, nil
	}
	return nil, errors.New("No Configuration loaded")
}

func LoadConfiguration(path string) (*Configuration, error) {

	if config != nil {
		return config, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	config = &Configuration{}
	err = decoder.Decode(config)

	if err != nil {
		return nil, err
	}

	return config, nil

}
