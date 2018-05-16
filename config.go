package main

import (
	"github.com/gpestana/envaws/aws"
	"log"
	"time"
)

type ConfigManager interface {
	GetConfigurations() (string, error)
	StartPolling()
}

type managerConfig struct {
	changedFn  func()
	pollingInt int
	accessKey  string
	secretKey  string
}

type S3Config struct {
	Env             map[string]string
	Interface       *aws.S3
	lastCheckedETag string
	config          managerConfig
}

func NewS3ConfigManager(interv int, bucket string, key string, aKey string, sKey string, f func()) S3Config {
	intf := aws.NewS3(bucket, key)

	return S3Config{
		Interface:       &intf,
		lastCheckedETag: "",
		config: managerConfig{
			changedFn:  f,
			accessKey:  aKey,
			secretKey:  sKey,
			pollingInt: interv,
		},
	}
}

func (c *S3Config) GetConfigurations() (string, error) {
	res, err := c.Interface.GetContent()
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func (c *S3Config) StartPolling() {
	// checks ETag every <c.interval> seconds
	for _ = range time.Tick(time.Duration(time.Duration(c.config.pollingInt) * time.Second)) {
		go func() {
			etag, err := c.Interface.GeCurrentETag()
			if err != nil {
				log.Println(err)
				return
			}

			if c.lastCheckedETag == "" {
				c.lastCheckedETag = etag
				return
			}

			if etag != c.lastCheckedETag {
				c.lastCheckedETag = etag
				c.config.changedFn()
			}
		}()
	}
}

type SSMConfig struct {
	Names     []string
	Interface *aws.SSM
	config    managerConfig
}

func NewSSMConfigManager(config C, f func()) SSMConfig {
	return SSMConfig{}
}

func (c *SSMConfig) GetConfigurations() (string, error) {
	return "", nil
}

func (c *SSMConfig) StartPolling() {}
