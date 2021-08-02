package main

import (
	"os"

	"github.com/manishmeganathan/weave/cmd"
)

func main() {
	defer os.Exit(0)
	cmd.Execute()
}
