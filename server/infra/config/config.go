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
	ListenAddress     string `yaml:"listenAddress"`
	ServerKeyFilepath string `yaml:"serverKeyFilepath"`
	ServerPemFilepath string `yaml:"serverPemFilepath"`
	ClientPemFilepath string `yaml:"clientPemFilepath"`
}

type MySQLConfig struct {
	SourceName      string `yaml:"sourceName"`
	MaxOpenConn     int    `yaml:"maxOpenConn"`
	MaxIdleConn     int    `yaml:"maxIdleConn"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type QuotaRepoConfig struct {
	UseMySQL   bool `yaml:"useMySQL"`
	FixedQuota int  `yaml:"fixedQuota"`
}

type StoreRepoConfig struct {
	UseRedis bool `yaml:"useRedis"`
}

type Config struct {
	ServerConfig    ServerConfig    `yaml:"server"`
	RedisConfig     RedisConfig     `yaml:"redis"`
	MySQLConfig     MySQLConfig     `yaml:"mysql"`
	QuotaRepoConfig QuotaRepoConfig `yaml:"quotaRepo"`
	StoreRepoConfig StoreRepoConfig `yaml:"storeRepo"`
}

func NewConfig() (c *Config) {
	err := lion.Load(file.NewSource(file.WithPath("config.yaml")))
	if err != nil {
		log.Fatal("Failed to load config. err:[%v]", err)
	}

	c = &Config{}

	err = lion.Get("server").Scan(&c.ServerConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	err = lion.Get("redis").Scan(&c.RedisConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	err = lion.Get("mysql").Scan(&c.MySQLConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	err = lion.Get("quotaRepo").Scan(&c.QuotaRepoConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	err = lion.Get("storeRepo").Scan(&c.StoreRepoConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	log.Infof("NewConfig success. config:[%+v]", c)

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
