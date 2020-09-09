package config

import (
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

var Map = new(ConfMap)

func Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	if err := yaml.NewDecoder(file).Decode(&Map); err != nil {
		return fmt.Errorf("could not decode config file: %v", err)
	}
	log.Printf("ConfigMap Loaded: [version: %s]", Map.App.Version)
	return nil
}
