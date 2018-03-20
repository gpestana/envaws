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

	// parses configurations for envaws
	//envawsConf := parseConf(cli.ConfPath)

	fn := func() {
		fmt.Println("envaws: Configurations changed, exiting process")
		os.Exit(0)
	}

	bucket, key, secretKey, accessKey, err := parseConfig(cli.ConfPath)
	if err != nil {
		log.Fatal(err)
	}

	config := NewConfig(5, bucket, key, secretKey, accessKey, fn)
	config.GetConfigurations()
	go config.StartPolling()

	// prepares new process
	cmd := exec.Command(cli.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// populate process env environment
	env := os.Environ()
	env = append(env, fmt.Sprintf("VAR1=%v", "var_1"))
	env = append(env, fmt.Sprintf("VAR2=%v", "var_2"))
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

	log.Println("Done")
}

func parseConfig(path string) (string, string, string, string, error) {
	type C struct {
		Bucket    string `json:"bucket"`
		Key       string `json:"key"`
		SecretKey string `json:"secret_key"`
		AccessKey string `json:"access_key"`
	}

	c := C{}
	js, err := ioutil.ReadFile(path)
	if err != nil {
		return c.Bucket, c.Key, c.SecretKey, c.AccessKey, err
	}

	err = json.Unmarshal(js, &c)
	if c.Bucket == "" {
		return c.Bucket, c.Key, c.SecretKey, c.AccessKey, errors.New("conf: bucket name should be provided")
	}

	if c.Key == "" {
		return c.Bucket, c.Key, c.SecretKey, c.AccessKey, errors.New("conf: key should be provided")
	}

	return c.Bucket, c.Key, c.SecretKey, c.AccessKey, nil
}
