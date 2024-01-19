package config

import (
	"log"
	"os"
)

type Conf struct {
	BaseUrl     string `yaml:"base_url"`
	DefaultLang string `yaml:"default_lang" default:"uz"`
	DBUsername  string `yaml:"db_username"`
	DBPassword  string `yaml:"db_password"`
	DBName      string `yaml:"db_name"`
	DBPort      string `yaml:"db_port"`
	Port        string `yaml:"port"`
	JWTKey      string `yaml:"jwt_key"`
}

func GetConf() *Conf {

	cfg := Conf{}

	yamlFile, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &cfg
}
