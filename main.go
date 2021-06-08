package main

import (
	"flag"
	"log"
	"net/url"
	"os"
)

func main() {
	// parse flags
	flagListen := flag.String("l", ":8000", "listening address")
	flagBackend := flag.String("b", "", "Target backend URL")
	flagScriptFile := flag.String("f", "", "Script file to run")
	flagRecordDir := flag.String("r", "", "Directory to store traffic records")
	flag.Parse()

	// override by environment variables
	listenAddr := os.Getenv("LISTEN")
	if listenAddr == "" {
		listenAddr = *flagListen
	}
	backend := os.Getenv("BACKEND")
	if backend == "" {
		backend = *flagBackend
	}

	backendURL, err := url.Parse(backend)
	if err != nil {
		log.Fatal("failed to parse backend URL: ", err)
	}

	scriptFilename := os.Getenv("SCRIPT_FILE")
	if scriptFilename == "" {
		scriptFilename = *flagScriptFile
	}

	recordDir := os.Getenv("RECORD_DIR")
	if recordDir == "" {
		recordDir = *flagRecordDir
	}

	StartServer(listenAddr, backendURL, scriptFilename, recordDir)

}
