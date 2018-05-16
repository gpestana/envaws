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

	config := NewConfig(c.PollingInterval, c.Bucket, c.Key, c.SecretKey, c.AccessKey, fn)

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
	PollingInterval int    `json:"polling_interval"`
	Bucket          string `json:"bucket"`
	Key             string `json:"key"`
	SecretKey       string `json:"secret_key"`
	AccessKey       string `json:"access_key"`
}

func parseConfig(path string) (C, error) {
	c := C{}
	js, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(js, &c)
	if c.Bucket == "" {
		return c, errors.New("conf: bucket name should be provided")
	}

	if c.Key == "" {
		return c, errors.New("conf: key should be provided")
	}

	if c.PollingInterval == 0 {
		return c, errors.New("conf: polling_interval should be provided")
	}

	return c, nil
}
