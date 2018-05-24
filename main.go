package main

import (
	"fmt"
	"github.com/gpestana/envaws/cli"
	"github.com/gpestana/envaws/client"
	"github.com/gpestana/envaws/config"
	"log"
	"os"
	"os/exec"
)

func main() {
	cli := cli.CLI{}
	err := cli.Run()
	if err != nil {
		os.Exit(1)
	}

	fn := func() {
		fmt.Println("envaws: Configurations changed, exiting process")
		os.Exit(0)
	}

	c, err := config.New(cli.ConfPath, cli.Service)
	if err != nil {
		log.Fatal(err)
	}

	aws, err := client.New(c, fn)
	if err != nil {
		log.Fatal(err)
	}

	confs, err := aws.GetConfigurations()
	if err != nil {
		log.Fatal(err)
	}

	go aws.StartPolling()

	cmd := exec.Command(cli.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	env := os.Environ()
	env = append(env, fmt.Sprintf("CONFIGS=%v", string(confs)))
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
