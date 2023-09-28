package main

import (
	"fmt"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("envdir: can't read envs: %v\n", err.Error())
		return
	}
	os.Exit(RunCmd(os.Args[2:], env))
}
