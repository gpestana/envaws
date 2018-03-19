package main

import (
	"github.com/gpestana/envaws/aws"
)

type Config struct {
	Env       map[string]string
	Interface *aws.S3
	changedFn func()
	accessKey string
	secretKey string
}

func NewConfig(bucket string, key string, aKey string, sKey string, f func()) Config {
	intf := aws.New(bucket, key)

	return Config{
		Interface: &intf,
		changedFn: f,
		accessKey: aKey,
		secretKey: sKey,
	}
}

func (c *Config) GetConfigurations() {}
func (c *Config) StartPolling() {
	c.changedFn()
}
