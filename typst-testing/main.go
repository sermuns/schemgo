package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	inp := bytes.NewBufferString("#set page(width: auto, height: auto, margin: 0cm, fill: none)\n")

	args := os.Args
	if len(args) < 2 {
		log.Fatal("No input specified.")
	}

	inp.WriteString(args[1])

	cmd := exec.Command("typst", "compile", "--format", "svg", "-", "out.svg")
	cmd.Stdin = inp
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	time.Sleep(time.Second)
	inp.WriteString(args[1])
	time.Sleep(time.Second)
	inp.WriteString(args[1])
	time.Sleep(time.Second)
	inp.WriteString(args[1])
	if err != nil {
		log.Fatal(err)
	}
}
