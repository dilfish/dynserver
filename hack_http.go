package main

import (
	"log"
	"net"
	"net/http"
)

func ConnState(conn net.Conn, state http.ConnState) {
	log.Println("connection state:", conn.RemoteAddr(), state)
}
