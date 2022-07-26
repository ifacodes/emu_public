package main

import (
	"io/ioutil"
	"log"
	"os"

	. "github.com/ifamakes/emu/pkg/hardware"
)

func main() {
	log_file, err := os.Create("emu_log")
	if err != nil {
		panic(err)
	}
	defer log_file.Close()
	log.SetFlags(0)
	log.SetOutput(log_file)

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	compare, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer compare.Close()

	g := NewGBC(file, compare)

	for {
		g.Step()
	}

}
