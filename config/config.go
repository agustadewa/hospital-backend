package config

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

var JWTSecretKey = os.Getenv("JWT_SECRET")

var CONFIG = ReadConfig()

const configFileName string = "config.json"

type TRepoConfig struct {
	DBName   string `json:"db_name"`
	CollName string `json:"coll_name"`
}

type TConfig struct {
	ServiceHost  string                 `json:"service_host"`
	ServicePort  string                 `json:"service_port"`
	Repositories map[string]TRepoConfig `json:"repositories"`
	MongoURL     string                 `json:"mongo_url"`
}

func ReadConfig() TConfig {
	file, err := os.Open(configFileName)
	if err != nil {
		logrus.Error(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	decode := json.NewDecoder(file)

	config := TConfig{}
	if err := decode.Decode(&config); err != nil {
		logrus.Error(err)
	}

	return config
}
