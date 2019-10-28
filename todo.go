package main

import (
	"flag"
	"log"
	"os"
)

var (
	mkfile = flag.String("mkfile", "", "Alternate Makefile to use")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(0)
	}
	cmd, err := newCommand(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	cmd.run()
}
