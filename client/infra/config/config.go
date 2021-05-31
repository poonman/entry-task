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

type Config struct {
	ServerConfig ServerConfig `yaml:"server" json:"server"`
	LogConfig LogConfig `yaml:"level" json:"level"`
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

type LogConfig struct {
	Level string `yaml:"level" json:"level"`
}
type ServerConfig struct {
	Address       string `yaml:"address" json:"address"`
	MaxActiveConn int    `yaml:"maxActiveConn" json:"maxActiveConn"`
	EnableTls     bool   `yaml:"enableTls" json:"enableTls"`
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

	err = lion.Get("log").Scan(&c.LogConfig)
	if err != nil {
		log.Errorf("Failed to scan config. err:[%v]", err)
		return
	}

	go c.watchLogConfig()

	log.Infof("NewConfig success... config:[%s]", c)

	return
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

func (c *Config) LoadTlsConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		panic("unable to load pem and key")
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
