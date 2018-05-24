package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	cli := CLI{}
	err := cli.Run()
	if err != nil {
		os.Exit(1)
	}

	fn := func() {
		fmt.Println("envaws: Configurations changed, exiting process")
		os.Exit(0)
	}

	c, err := parseConfig(cli.ConfPath)
	if err != nil {
		log.Fatal(err)
	}

	var config ConfigManager
	// creates configuration manager based on type of service
	switch cli.Service {
	case "s3":
		config = NewS3ConfigManager(c, fn)
	case "ssm", "":
		// default service is ssm
		config = NewSSMConfigManager(c, fn)
	default:
		log.Fatal(
			fmt.Sprintf("Service %v not supported. Use 'ssm' or 's3' instead", cli.Service))
	}

	conf, err := config.GetConfigurations()
	if err != nil {
		log.Fatal(err)
	}
	go config.StartPolling()

	cmd := exec.Command(cli.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	env := os.Environ()
	env = append(env, fmt.Sprintf("CONFIGS=%v", conf))
	cmd.Env = env

	// starts new process
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

type C struct {
	// configuration for accessing parameters in S3
	S3 struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	} `json:"s3"`

	// configuration for accessing parameters in S3
	Ssm struct {
		Names []string `json:"parameters_name"`
	} `json:"parameters_name"`

	// AWS-wide configurations
	SecretKey       string `json:"secret_key"`
	AccessKey       string `json:"access_key"`
	PollingInterval int    `json:"polling_interval"`
	Region          string `json:"region"`
}

func parseConfig(path string) (C, error) {
	c := C{}
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
