package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	LogLevel   string           `json:"log_level"`
	HttpServer HttpServerConfig `json:"http_server"`
	Database   DatabaseConfig   `json:"database"`
	Postgres   PostgresConfig   `json:"postgres"`
	Mongo      MongoConfig      `json:"mongo"`
}

func NewConfig() *Config {
	conf := new(Config)

	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

type HttpServerConfig struct {
	ListenAddress string `json:"listen_address"`
}

type DatabaseConfig struct {
	Name string `json:"name"`
}

type PostgresConfig struct {
	URI string `json:"URI"`
}

type MongoConfig struct {
	URI string `json:"URI"`
	DB  string `env:"DB"`
}
