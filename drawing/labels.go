package drawing

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

var (
	preamble = bytes.NewBufferString(`
#set page(width: auto, height: auto, margin: 0cm, fill: none)
		`)
)

func renderTypstSvg(body []byte) []byte {
	inp := bytes.NewBuffer(preamble.Bytes())

	args := os.Args
	if len(args) < 2 {
		log.Fatal("No input specified.")
	}

	inp.Write(body)

	cmd := exec.Command("typst", "compile", "--format", "svg", "-", "out.svg")

	var outBuf, errBuf bytes.Buffer
	cmd.Stdin = inp
	cmd.Stdout = &outBuf // Redirect Stdout to outBuf
	cmd.Stderr = &errBuf // Redirect Stderr to errBuf

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed: %s\nStderr: %s", err, errBuf.String())
	}

	return outBuf.Bytes()
}
