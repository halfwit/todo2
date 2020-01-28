package main

import (
	"log"
	"os"
)

func main() {
	cmd, err := newCommand(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(cmd.run())
}
