// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var (
	serverRoot       string // Absolute path of server root.
	staticFolderPath string // Absolute path of static file folder.
	faviconPath      string // Absolute path of "favicon.ico".
	indexTmplPath    string // Absolute path of index HTML template file.
	addr             = flag.String("addr", ":80", "http service address")
)

func init() {
	// Get absolute path of server root(current executable).
	serverRoot, _ = GetCurrentExecDir()
	// Get static folder path.
	staticFolderPath = path.Join(serverRoot, "./dist/spa")
	// Get favicon.ico path.
	faviconPath = path.Join(serverRoot, "favicon.ico")
	// Get index template file path.
	indexTmplPath = path.Join(serverRoot, "./templates/index.tmpl")
}

func GetCurrentExecDir() (dir string, err error) {
	p, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	dir = filepath.Dir(absPath)
	return dir, nil
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	SetUp()
	flag.Parse()
	//	http.HandleFunc("/", serveHome)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(staticFolderPath))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
