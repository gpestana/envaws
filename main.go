package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	cli := CLI{}
	cli.Run()

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
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done")
}
