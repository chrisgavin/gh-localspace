package main

import (
	"github.com/chrisgavin/gh-localspace/cmd"
)

func main() {
	rootCommand, err := cmd.NewRootCommand()
	if err != nil {
		panic(err)
	}
	rootCommand.Run()
}
