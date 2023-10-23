package main

import (
	"log"
	"net/http"
	"strconv"
)

func StatusCode(w http.ResponseWriter, r *http.Request) {
	if len(r.RequestURI) <= len("/status/") {
		log.Println("bad status:", r.RequestURI)
		w.Write([]byte("bad status"))
		return
	}
	codeStr := r.RequestURI[len("/status/"):]
	code, err := strconv.ParseUint(codeStr, 10, 32)
	if err != nil {
		log.Println("bad code:" + r.RequestURI)
		return
	}
	if code > 1000 {
		code = 1000
	}
	w.WriteHeader(int(code))
	w.Write([]byte("ok"))
}
