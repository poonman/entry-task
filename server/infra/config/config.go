package config

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/poonman/entry-task/dora/misc/lion"
	"github.com/poonman/entry-task/dora/misc/lion/source/file"
	"github.com/poonman/entry-task/dora/misc/log"
	"io/ioutil"
)

type ServerConfig struct {
	ListenAddress string `yaml:"listenAddress" json:"listenAddress"`
	EnableTls     bool   `yaml:"enableTls" json:"enableTls"`
}

type LogConfig struct {
	Level string `yaml:"level" json:"level"`
}

type MySQLConfig struct {
	SourceName      string `yaml:"sourceName" json:"sourceName"`
	MaxOpenConn     int    `yaml:"maxOpenConn" json:"maxOpenConn"`
	MaxIdleConn     int    `yaml:"maxIdleConn" json:"maxIdleConn"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime"`
}

type RedisConfig struct {
	Address   string `yaml:"address" json:"address"`
	Password  string `yaml:"password" json:"password"`
	MaxIdle   int    `yaml:"maxIdle" json:"maxIdle"`
	MaxActive int    `yaml:"maxActive" json:"maxActive"`
}

type QuotaRepoConfig struct {
	UseMySQL   bool `yaml:"useMySQL" json:"useMySQL"`
	FixedQuota int  `yaml:"fixedQuota" json:"fixedQuota"`
}

type StoreRepoConfig struct {
	UseRedis bool `yaml:"useRedis" json:"useRedis"`
}

type RateLimiterRepoConfig struct {
	Capacity int `yaml:"capacity" json:"capacity"`
}

type Config struct {
	ServerConfig          ServerConfig          `yaml:"server" json:"server"`
	LogConfig             LogConfig             `yaml:"log" json:"log"`
	RedisConfig           RedisConfig           `yaml:"redis" json:"redis"`
	MySQLConfig           MySQLConfig           `yaml:"mysql" json:"mysql"`
	QuotaRepoConfig       QuotaRepoConfig       `yaml:"quotaRepo" json:"quotaRepo"`
	StoreRepoConfig       StoreRepoConfig       `yaml:"storeRepo" json:"storeRepo"`
	RateLimiterRepoConfig RateLimiterRepoConfig `yaml:"rateLimiter" json:"rateLimiter"`
}

func (c *Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return ""
	}
	return out.String()
}

func NewConfig() (c *Config) {
	err := lion.Load(file.NewSource(file.WithPath("config.yaml")))
	if err != nil {
		log.Fatal("Failed to load config. err:[%v]", err)
	}

	c = &Config{}

	err = lion.Get("log").Scan(&c.LogConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

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

	err = lion.Get("rateLimiter").Scan(&c.RateLimiterRepoConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	go c.watchLogConfig()

	log.Infof("NewConfig success. config:[%+v]", c)

	return c
}

func (c *Config) watchLogConfig() {
	w, err := lion.Watch("log")
	if err != nil {
		log.Errorf("Failed to watch log config. err:[%v]", err)
		return
	}

	for {
		value, err := w.Next()
		if err != nil {
			log.Errorf("Failed to watch next. err:[%v]", err)
			continue
		}

		err = value.Scan(&c.LogConfig)
		if err != nil {
			log.Errorf("Failed to scan config. err:[%v]", err)
			continue
		}

		log.SetLevelByString(c.LogConfig.Level)
	}
}

func (c *Config) LoadTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Fatal("Failed to read server.pem and server.key. err:[%v]", err)
	}
	certBytes, err := ioutil.ReadFile("client.pem")
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
