package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var snippets []string
	err = json.Unmarshal(data, &snippets)
	if err != nil {
		log.Fatal(err)
	}

	var cmdIn strings.Builder
	cmdIn.WriteString("#set page(width:auto,height:auto,margin:0cm,fill:none)\n")

	for _, snippet := range snippets {
		cmdIn.WriteString(snippet + "\n#pagebreak()\n")
	}

	os.Mkdir("/dev/shm/schemgo", 0755)
	cmd := exec.Command("typst", "compile", "--format", "svg", "-", "/dev/shm/schemgo/schemgo{0p}.svg")
	cmd.Stdin = strings.NewReader(cmdIn.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Took", time.Since(start))
}
