package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Port            string  `yaml:"port"`
	MongoDb         MongoDb `yaml:"mongoDb"`
	LogPath         string  `yaml:"logPath"`
	ArticleListURL  string  `yaml:"articleListURL"`
	ArticleURL      string  `yaml:"articleURL"`
	ArticleInterval int     `yaml:"articleInterval"`
}

type MongoDb struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	DbName     string `yaml:"dbName"`
}

var Conf *Config

func GetConfig(configFile string) *Config {
	Conf = &Config{}
	if configFile != "" {
		Conf.GetConfFromFile(configFile)
	}

	logFilePath := Conf.LogPath

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
		log.Panic("Failed to log to file, using default stderr", err)
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	return Conf
}

func (c *Config) GetConfFromFile(configFile string) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}
}
