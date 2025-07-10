package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	SIS struct {
		Host    string `toml:"host"`
		Port    int    `toml:"port"`
		Service string `toml:"service"`
	} `toml:"sis"`
	RecSIS struct {
		Host   string `toml:"host"`
		Port   int    `toml:"port"`
		DBName string `toml:"dbname"`
	} `toml:"recsis"`
	MeiliSearch struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"meilisearch"`
}

func loadConfig() (config, error) {
	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()
	if len(*configPath) == 0 {
		log.Fatal("Config file path is required")
	}

	var conf config
	_, err := toml.DecodeFile(*configPath, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to load config file: %w", err)
	}
	return conf, nil
}
