package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"net/http"
	"path"
	"strconv"

	handlers "cornchip.com/pkmneventserver/Handlers"
)

const DEFAULT_HOST string = "127.0.0.12"
const DEFAULT_PORT int = 8080
const DEFAULT_LOG_LEVEL string = "error"
const DEFAULT_HUMAN_LOGS  bool = true
const DEFAULT_CARD_ROOT string = "./cards"

func main() {
	var addressFlag string
	var portFlag int
	var logLevel string
	var humanLogs bool
	var cardRoot string

	flag.StringVar(&addressFlag, "addr", DEFAULT_HOST, "address to listen on")
	flag.IntVar(&portFlag, "port", DEFAULT_PORT, "port to listen on")
	flag.StringVar(&logLevel, "log", DEFAULT_LOG_LEVEL, "log level")
	flag.BoolVar(&humanLogs, "readable", DEFAULT_HUMAN_LOGS, "use human readable logs")
	flag.StringVar(&cardRoot, "cardroot", DEFAULT_CARD_ROOT, "root path to WonderCards")

	flag.Parse()

	// Setup log config
	handlers.SetHumanLogs(humanLogs)
	handlers.SetLogLevel(logLevel)

	// Check if directories exist
	handlers.SetRootPath( cardRoot )

	rootPath := handlers.GetRootPath()
	f, err := os.Open( path.Join(rootPath, "geniv") )
	// folder doesnt exist, make it
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating Gen IV card directory...")
		err = os.MkdirAll( rootPath, 0755 )
		if err != nil {
			errorMessage := "Unable to automatically create directory " +
							rootPath + ". Error: " + err.Error()
			log.Fatal(errorMessage)
		}
	} else if err != nil {
		errorMessage := "Unspecified error: " + err.Error()
		log.Fatal(errorMessage)
	}

	f.Close()

	listenAddress := addressFlag + ":" + strconv.Itoa(portFlag)

	http.HandleFunc("/nas/ac", handlers.HandleNAS)
	http.HandleFunc("/runConnTest", handlers.HandleConnTest)
	http.HandleFunc("/dls1/download", handlers.HandleDls1)
	http.HandleFunc("/gs2", handlers.HandleGTS)

	http.HandleFunc("/", handlers.HandleOther)

	fmt.Println("Running fake event server on " + listenAddress + "...")

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}