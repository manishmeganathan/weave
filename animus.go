package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/manishmeganathan/animus/cli"
)

func main() {
	defer os.Exit(0)
	rand.Seed(time.Now().UTC().UnixNano())
	log.SetPrefix("animus ")

	cli.Execute()
}
