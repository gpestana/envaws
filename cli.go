package main

import (
	"errors"
	"flag"
)

type CLI struct {
	ConfPath string
	Command  string
	Service  string
}

func (cli *CLI) Run() error {
	conf := flag.String("conf", "", "Path for envaws configurations [required]")
	command := flag.String("command", "", "Command to be called [required]")
	serv := flag.String("service", "", "Ssrvice to use (ssm or s3); default: ssm")

	flag.Parse()

	cli.ConfPath = *conf
	cli.Command = *command
	cli.Service = *serv

	if cli.ConfPath == "" {
		flag.PrintDefaults()
		return errors.New("-conf flag missing")
	}

	if cli.Command == "" {
		flag.PrintDefaults()
		return errors.New("-command flag missing")
	}
	return nil
}
