package config

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/misc/lion"
	"github.com/poonman/entry-task/dora/misc/lion/source/file"
	"io/ioutil"
)

type ServerConfig struct {
	ListenAddress string `yaml:"listenAddress"`
	ServerKeyFilepath string `yaml:"serverKeyFilepath"`
	ServerPemFilepath string `yaml:"serverPemFilepath"`
	ClientPemFilepath string `yaml:"clientPemFilepath"`
}

type MySQLConfig struct {
	SourceName string `yaml:"sourceName"`
}

type RedisConfig struct {
	Address string `yaml:"address"`
	Password string `yaml:"password"`
}

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
	RedisConfig RedisConfig `yaml:"redis"`
	MySQLConfig MySQLConfig `yaml:"mysql"`
}

func NewConfig() (c *Config){
	err := lion.Load(file.NewSource(file.WithPath("config.yaml")))
	if err != nil {
		log.Fatal("Failed to load config. err:[%v]", err)
	}

	c = &Config{}

	err = lion.Get().Scan(c)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	return c
}

func (c *Config) LoadTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(c.ServerConfig.ServerPemFilepath, c.ServerConfig.ServerKeyFilepath)
	if err != nil {
		log.Fatal("Failed to read server.pem and server.key. err:[%v]", err)
	}
	certBytes, err := ioutil.ReadFile(c.ServerConfig.ClientPemFilepath)
	if err != nil {
		log.Fatal("Failed to read client.pem, err:[%v]", err)
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("Failed to parse root certificate.")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}

	return config
}