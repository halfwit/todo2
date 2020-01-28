package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()

	cmd, err := newCommand(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	s, err := cmd.run()
	if err != nil {
		log.Fatal(err)
	}

	// TODO(halfwit) Switch to our own writer type and pass in os.Stdout
	fmt.Printf("%s", s)
}
