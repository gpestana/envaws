package main

import (
	"github.com/gpestana/envaws/aws"
	"log"
	"time"
)

type Config struct {
	Env             map[string]string
	Interface       *aws.S3
	changedFn       func()
	accessKey       string
	secretKey       string
	lastCheckedETag string
	pollingInt      int
}

func NewConfig(interv int, bucket string, key string, aKey string, sKey string, f func()) Config {
	intf := aws.New(bucket, key)

	return Config{
		Interface:       &intf,
		changedFn:       f,
		accessKey:       aKey,
		secretKey:       sKey,
		pollingInt:      interv,
		lastCheckedETag: "",
	}
}

func (c *Config) GetConfigurations() {
	res, err := c.Interface.GetContent()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.String())
}

func (c *Config) StartPolling() {
	// checks ETag every <c.interval> seconds
	for _ = range time.Tick(time.Duration(time.Duration(c.pollingInt) * time.Second)) {
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
				c.changedFn()
			}
		}()
	}
}
