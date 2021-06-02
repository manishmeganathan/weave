package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/manishmeganathan/animus/cmd"
)

func main() {
	defer os.Exit(0)
	rand.Seed(time.Now().UTC().UnixNano())
	log.SetPrefix("animus ")

	cmd.Execute()
}
