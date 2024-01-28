package configs

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Port                   string `yaml:"port"`
	Address                string `yaml:"address"`
	TokenKey               string `yaml:"token-key"`
	DatabasePort           string `yaml:"database-port"`
	DatabaseAddress        string `yaml:"database-address"`
	DatabaseUser           string `yaml:"database-user"`
	DatabasePassword       string `yaml:"database-password"`
	DatabaseName           string `yaml:"database-name"`
	StorageServiceID       string `yaml:"storage-service-id"`
	StorageServiceSecret   string `yaml:"storage-service-secret"`
	StorageServiceEndpoint string `yaml:"storage-service-endpoint"`
	StorageServiceBucket   string `yaml:"storage-service-bucket"`
}

var Config = new(Configuration)

func ParseConfig() error {
	yamlFile, err := ioutil.ReadFile("./configs/config.yaml")
	if err != nil {
		logrus.Printf("Can not read configuration file")
		return err
	}
	err = yaml.Unmarshal(yamlFile, Config)
	if err != nil {
		logrus.Println("Can not unmarshal configuration file")
		return err
	}

	return nil
}
