package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	if err := execute(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
