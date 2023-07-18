// http just redirect requests to https
package main

import (
	"log"
	"net"
	"net/http"
)

type RedirectHandler struct {
	Same http.Handler
}

func (rd *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("http request is:", r.RemoteAddr, "->", r.Host, r.RequestURI)
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host == "127.0.0.1" {
		rd.Same.ServeHTTP(w, r)
		return
	}
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, 302)
}

func Redirect(h http.Handler) {
	var rd RedirectHandler
	rd.Same = h
	err := http.ListenAndServe(":1080", &rd)
	if err != nil {
		panic("listen and serve http error: " + err.Error())
	}
}
