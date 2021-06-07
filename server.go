package main

import (
	"fmt"
	"github.com/lechuckroh/websocket-proxy-server/session"
	esbuild2 "github.com/lechuckroh/websocket-proxy-server/session/esbuild"
	"log"
	"net/http"
	"net/url"
)

type sessionHandler struct {
	BackendURL     *url.URL
	ScriptFilename string
}

func (h *sessionHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	code, err := h.compileScript()
	if err != nil {
		log.Printf("[ERR] failed to compile %s: %v", h.ScriptFilename, err)
		return
	}

	s, err := session.NewSession(h.BackendURL, code, rw, req)
	if err != nil {
		log.Printf("[ERR] failed to create session: %v", err)
		return
	}

	s.Start()
}

func (h *sessionHandler) compileScript() (string, error) {
	if h.ScriptFilename == "" {
		return "", nil
	}

	compiler := esbuild2.CompilerImpl{}
	code, buildResult := compiler.Compile(h.ScriptFilename)

	errorCount := len(buildResult.Errors)
	if errorCount > 0 {
		for _, msg := range buildResult.Errors {
			log.Printf("[ERR] %v", msg.Text)
		}
		return "", fmt.Errorf("%d compile errors", errorCount)
	}

	warningCount := len(buildResult.Warnings)
	if warningCount > 0 {
		for _, msg := range buildResult.Warnings {
			log.Printf("[WARN] %v", msg.Text)
		}
	}
	return code, nil
}

// StartServer starts a websocket server.
func StartServer(
	listenAddr string,
	backendURL *url.URL,
	scriptFilename string,
	recordDir string,
) {
	log.Printf("starting server on %s", listenAddr)

	handler := &sessionHandler{
		BackendURL:     backendURL,
		ScriptFilename: scriptFilename,
	}
	if err := http.ListenAndServe(listenAddr, handler); err != nil {
		log.Fatal(err)
	}
}
