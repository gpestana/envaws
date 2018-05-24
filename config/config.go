package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type C struct {
	// configuration for accessing parameters in S3
	S3 struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	} `json:"s3"`

	// configuration for accessing parameters in S3
	Ssm struct {
		Parameters []string `json:"parameters"`
	} `json:"ssm"`

	// AWS-wide configurations
	SecretKey       string `json:"secret_key"`
	AccessKey       string `json:"access_key"`
	PollingInterval int    `json:"polling_interval"`
	Region          string `json:"region"`
	Service         string
}

func (c *C) setService(s string) {
	c.Service = s
}

func New(path string, srv string) (C, error) {
	c := C{}
	c.Service = srv

	js, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(js, &c)
	if c.S3.Bucket == "" {
		return c, errors.New("conf: bucket name must be provided")
	}

	if c.S3.Key == "" {
		return c, errors.New("conf: key must be provided")
	}

	if c.PollingInterval == 0 {
		return c, errors.New("conf: polling_interval must be provided")
	}

	return c, nil
}
