package main

import (
	"flag"
)

type CLI struct {
	Conf    string
	Command string
}

func (cli *CLI) Run() {
	conf := flag.String("conf", "", "configurations")
	command := flag.String("command", "", "command to run") //TODO: funny name, change
	flag.Parse()

	cli.Conf = *conf
	cli.Command = *command

}
