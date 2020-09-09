package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type ConfMap struct {
	App struct {
		Version string `yaml:"version"`
		Debug   bool   `yaml:"debug"`
		Env     string `yaml:"env"`
	} `yaml:"app"`
	Secrets struct {
		Db struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
			Name string `yaml:"name"`
		} `yaml:"db"`
		Oauth struct {
			Discord string `yaml:"discord"`
			Google  string `yaml:"google"`
			Spotify string `yaml:"spotify"`
		} `yaml:"oauth"`
		Redis struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
			Pass string `yaml:"pass"`
		} `yaml:"redis"`
		JWT struct {
			ExpireTime            int    `yaml:"expire_time"`
			RefreshTokenValidTime int    `yaml:"refresh_token_valid_time"`
			PrivateKeyPath        string `yaml:"private_key_path"`
			PublicKeyPath         string `yaml:"public_key_path"`
		} `yaml:"jwt"`
		SentryDsn      string `yaml:"sentry_dsn"`
		HcaptchaSecret string `yaml:"hcaptcha_secret"`
	} `yaml:"secrets"`
	StoragePath string `yaml:"storage_path"`
}

var Map *ConfMap

func init() {
	configFileName := flag.String("config-file", "config.yml", "config.yaml file")
	flag.Parse()
	file, err := os.Open(*configFileName)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not open config.yml file : %v", err))
	}
	if err := yaml.NewDecoder(file).Decode(&Map); err != nil {
		log.Fatal(fmt.Sprintf("could not decode config file : %v", err))
	}
	log.Printf("ConfigMap Loaded: [version: %s]", Map.App.Version)
}
