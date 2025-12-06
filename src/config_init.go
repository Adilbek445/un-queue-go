package main

import (
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}

	Queue struct {
		DataDir string
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	_, errInt := strconv.Atoi(cfg.Server.Port)

	if errInt != nil {
		panic(errInt)
	}

	if !isDirExists(cfg.Queue.DataDir) {
		panic("The directory specified in queue.dataDir does not exist.")
	}

	return &cfg, nil
}
