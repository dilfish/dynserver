package main

import (
	dnet "github.com/dilfish/tools/net"
	"io"
	"log"
	"net/http"
	"strings"
)

type HttpsHandler struct {
	u *dnet.UploaderService
	C *MongoClient
}

func (h *HttpsHandler) Uploader(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, dnet.GetUploadPage("上传文件", "/upload"))
		return
	}
	h.u.Handler(w, r)
}

func (h *HttpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request is: ", r.RemoteAddr, "->", r.Host, r.RequestURI)
	log.Println("content length is:", r.ContentLength)
	log.Println("header is: ", r.Header)
	if r.Host != "ls.dev.ug" && r.Host != "ls4.dev.ug" {
		// return
	}
	if strings.Index(r.RequestURI, "/memfile/") == 0 {
		MemFileHandler(w, r)
		return
	}
	if strings.Index(r.RequestURI, "/upload") == 0 {
		h.Uploader(w, r)
		return
	}
	if strings.Index(r.RequestURI, "/t/") == 0 {
		h.Msg(w, r)
		return
	}
	if r.RequestURI == "/t" {
		h.Msg(w, r)
		return
	}
	if r.RequestURI == "/ip" {
		CFIPHandler(w, r)
		return
	}
	fs := http.FileServer(http.Dir("/root/go/src/dynserver"))
	fs.ServeHTTP(w, r)
}
