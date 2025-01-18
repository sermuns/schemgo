package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

const (
	defaultAdress = "localhost:3000"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve circuit diagram as html on local network, rebuild and live-reload on change. Default: http://" + defaultAdress,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inFilePath = args[0]

		go watchInput(inFilePath)

		http.HandleFunc("/", serveSchematic)
		http.HandleFunc("/events", reloadHandler)

		address := cmd.Flag("bind").Value.String()

		go func() {
			time.Sleep(time.Second)
			if cmd.Flag("open").Value.String() == "true" {
				openBrowser("http://" + address)
			}
		}()

		log.Fatal(http.ListenAndServe(address, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("bind", "b", defaultAdress, "Bind to this address")
	serveCmd.Flags().BoolP("open", "o", false, "Open browser")
}

var (
	connections = &Connections{
		clients: make(map[chan []byte]struct{}),
	}

	inFilePath      = ""
	latestInContent []byte

	latestSchematic []byte
)

func updateSchematic() {
	inContent, err := os.ReadFile(inFilePath)
	if err != nil {
		log.Println("Error reading input file:", err)
		return
	}

	if bytes.Equal(inContent, latestInContent) {
		log.Println("content equal, dont care")
		return
	}

	log.Println("Updating schematic")

	latestInContent = inContent
	latestSchematic = writeSchematic(inContent)

	connections.broadcast(latestSchematic)
}

func serveSchematic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Write([]byte(`
<!DOCTYPE html>
<script>
function connect(){
	const eventSource = new EventSource('http://` + defaultAdress + `/events')	
	eventSource.onmessage = (event) => {
		document.querySelector('svg').outerHTML = event.data
	}
	eventSource.onerror = (err) => {
		console.error(err)
	}
}

connect()
</script>
`))
	w.Write(latestSchematic)
}

type Connections struct {
	clients map[chan []byte]struct{}
	mutex   sync.Mutex
}

func (c *Connections) add(client chan []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.clients[client] = struct{}{}
}

func (c *Connections) remove(client chan []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.clients, client)
	close(client)
}

func (c *Connections) broadcast(data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for client := range c.clients {
		select {
		case client <- data:
			// sending worked!
		default:
			// blocked or closed
			log.Println("deleting client")
			delete(c.clients, client)
			close(client)
		}
	}
}

func reloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	clientChan := make(chan []byte)
	connections.add(clientChan)
	defer connections.remove(clientChan)

	for data := range clientChan {
		fmt.Fprintf(w, "data: %s\n\n", data)
		w.(http.Flusher).Flush()
	}
}

func watchInput(inFilePath string) {
	updateSchematic()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer func() {
		log.Println("Closing watcher")
		watcher.Close()
	}()

	err = watcher.Add(filepath.Dir(inFilePath))
	if err != nil {
		log.Fatalf("Failed to watch input file: %v", err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Name != inFilePath {
				continue
			}

			if !event.Has(fsnotify.Write) {
				continue
			}

			updateSchematic()

		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
