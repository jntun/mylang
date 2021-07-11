package main

import (
	"fmt"
	jlang "github.com/jntun/mylang/lang"
	"io/ioutil"
	"log"
	"net/http"
)

func httpServer() {
	http.HandleFunc("/jlang", jlangHandler)
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

func jlangHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: the rest of the owl (build jlang interpreter hooking)
	switch r.Method {
	case "POST":
		runScript(w, r)
	case "PUT":
		runScript(w, r)
	case "GET":
		htmlFile, err := readHTML("jlang")
		if err != nil {
			log.Printf("Unable to open html file 'jlang.html': %s.\n", err)
			return
		}

		_, err = w.Write(htmlFile)

		if err != nil {
			log.Printf("Failed to write GET response: %s.\n", err)
			return
		}
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

func readHTML(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile("./public/" + filename + ".html")
	if err != nil {
		return nil, err
	}
	return data, nil
}
