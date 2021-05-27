package config

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/misc/lion"
	"github.com/poonman/entry-task/dora/misc/lion/source/file"
	"io/ioutil"
)

type Config struct {
	ServerConfig    ServerConfig    `yaml:"server"`
	BenchmarkConfig BenchmarkConfig `yaml:"benchmark"`
	CmdConfig       CmdConfig       `yaml:"cmd"`
}

type ServerConfig struct {
	Address       string `yaml:"address"`
	MaxActiveConn int    `yaml:"maxActiveConn"`
	EnableTls bool `yaml:"enableTls"`
}

type BenchmarkConfig struct {
	Concurrency              int `yaml:"concurrency"`
	RequestNumPerConcurrency int `yaml:"requestNumPerConcurrency"`
}

type CmdConfig struct {
	Commands []string `yaml:"commands"`
	Username string   `yaml:"username"`
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

	err = lion.Get("cmd").Scan(&c.CmdConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	log.Infof("NewConfig success... config:[%+v]", c)

	return
}

func (c *Config) LoadTlsConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		panic("unable to load pem and key" )
	}
	certBytes, err := ioutil.ReadFile("client.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	conf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	return conf
}