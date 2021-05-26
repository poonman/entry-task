package config

import (
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/misc/lion"
	"github.com/poonman/entry-task/dora/misc/lion/source/file"
)

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
	BenchmarkConfig BenchmarkConfig `yaml:"benchmark"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type BenchmarkConfig struct {
	Concurrency int `yaml:"concurrency"`
	RequestNumPerConcurrency int `yaml:"requestNumPerConcurrency"`
}

func NewConfig() (c *Config) {
	log.Infof("NewConfig begin...")

	err := lion.Load(file.NewSource(file.WithPath("config.yaml")))
	if err != nil {
		log.Fatal("Failed to load config. err:[%v]", err)
	}

	c = &Config{}

	err = lion.Get("server").Scan(&c.ServerConfig)
	//err = lion.Get().Scan(c)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	err = lion.Get("benchmark").Scan(&c.BenchmarkConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	log.Infof("NewConfig success... config:[%+v]", c)

	return
}