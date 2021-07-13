package main

import (
	"fmt"
	jlang "github.com/jntun/mylang/lang"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func httpServer() {
	http.HandleFunc("/jlang", jlangHandler)
	http.HandleFunc("/public/", publicHandler)
	http.HandleFunc("/", homeHandler)

	s := &http.Server{
		Addr: ":80",
	}

	log.Fatal(s.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var ret int
	var err error

	if r.Host == "localhost" {
		// DEBUG response
		ret, err = w.Write([]byte("Hello master..."))
	} else {
		// Public response
		ret, err = w.Write([]byte("Hello world from jntun.com..."))
	}

	if err != nil {
		log.Printf("Failure: %s.\n", err)
	}
	log.Printf("home: %d | %v\n", ret, r)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	fileURL := r.URL.String()[len("/public/"):]
	file := strings.Split(fileURL, ".")
	if len(file) < 1 {
		log.Printf("Invalid static file: %s.\n", r.URL.String())
		return
	}

	switch fileURL {
	case "ace.js":
		src := readStatic(file[0], file[1])
		w.Header().Set("Content-Type", "text/javascript")
		write(src, w)
	case "style.css":
		src := readStatic(file[0], file[1])
		w.Header().Set("Content-Type", "text/css")
		write(src, w)
	}
}

func jlangHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: the rest of the owl (build jlang interpreter hooking)
	switch r.Method {
	case "POST":
		runScript(w, r)
	case "PUT":
		runScript(w, r)
	case "GET":
		src := readStatic("jlang", "html")
		write(src, w)
	}
}

func runScript(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("Couldn't read body: %s.\n", err)
		return
	}
	log.Println("/jlang:", string(data))
	intptr := jlang.NewInterpreter()

	err = intptr.HookLogOut(w)
	if err != nil {
		log.Printf("Could not hook response writer: %s.\n", err)
		return
	}

	err = intptr.Interpret(string(data))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		return
	}
}

func readStatic(filename string, ext string) []byte {
	data, err := ioutil.ReadFile("./public/" + filename + "." + ext)
	if err != nil {
		log.Printf("Failed to load resource: %s.\n", err)
		return []byte("Failure to load resource.")
	}
	return data
}

func write(data []byte, w http.ResponseWriter) {
	if _, err := w.Write(data); err != nil {
		log.Printf("Failure writing response: %s.\n", err)
	}
}
