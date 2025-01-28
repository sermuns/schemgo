package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	var inp strings.Builder
	inp.WriteString("#set page(width: auto, height: auto, fill: none)\n")

	args := os.Args
	if len(args) < 2 {
		log.Fatal("No input specified.")
	}

	inp.WriteString(args[1])

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := exec.Command("typst", "compile", "--format", "svg", "-", "-")
			cmd.Stdin = strings.NewReader(inp.String())
			// cmd.Stdout = os.Stdout
			// cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	wg.Wait()

	fmt.Println("Took", time.Since(start))
}
