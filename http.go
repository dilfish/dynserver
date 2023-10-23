// http just redirect requests to https
package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var ConcurrentCount int
var ConcurrentLock sync.Mutex

const MaxConcurrentCount = 4

type RedirectHandler struct {
	Same http.Handler
}

func (rd *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ConcurrentLock.Lock()
	defer ConcurrentLock.Unlock()
	if ConcurrentCount > MaxConcurrentCount {
		time.Sleep(time.Second * 3)
		return
	}

	ConcurrentCount = ConcurrentCount + 1
	defer func() { ConcurrentCount = ConcurrentCount - 1 }()

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
	addr := ":80"
	if *FlagTestMode {
		addr = ":11080"
	}
	err := http.ListenAndServe(addr, &rd)
	if err != nil {
		panic("listen and serve http error: " + err.Error())
	}
}
