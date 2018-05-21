package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestParseS3OKConfig(t *testing.T) {
	p1 := "./conf.json"
	c1 := `
{
	"polling_interval": 3,
  "s3": {
	  "bucket": "gpestana-conf",
	  "key": "foo/env.json"
  },
  "ssm": {}
}
`
	expInt := 3
	expBucket := "gpestana-conf"
	expKey := "foo/env.json"

	defer os.Remove(p1)
	createConfFile(p1, c1)
	conf, err := parseConfig(p1)
	if err != nil {
		log.Fatal(err)
	}
	if conf.PollingInterval != expInt {
		t.Error(fmt.Sprintf("expected interval %v, got %v", expInt, conf.PollingInterval))
	}
	if conf.S3.Bucket != expBucket {
		t.Error(fmt.Sprintf("expected S3 bucket %v, got %v", expBucket, conf.S3.Bucket))
	}
	if conf.S3.Key != expKey {
		t.Error(fmt.Sprintf("expected S3 key %v, got %v", expKey, conf.S3.Key))
	}
}

// creates configuration file for testing purposes
func createConfFile(p string, s string) {
	file, err := os.Create(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(file, s)
}
