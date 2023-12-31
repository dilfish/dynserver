package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	dnet "github.com/dilfish/tools/net"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var SniList []string

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
	start := time.Now()
	defer func() {
		elaps := time.Now().Sub(start)
		us := elaps.Microseconds()
		requestDurationUs.Observe(float64(us))
	}()
	w.Header().Set("X-Server", "iPhone Xs")
	log.Println("request is: ", r.Method, r.RemoteAddr, "->", r.Host, r.RequestURI, r.ContentLength)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("split ip port error:", r.RemoteAddr, err)
		return
	}
	if r.Header.Get("X-No-Check-IP") != "true" {
		good := IsGoodIP(ip, r.Header)
		if !good {
			log.Println("bad ip:", ip)
			return
		}
	}
	log.Println("header is: ", r.RemoteAddr, r.Header)
	if r.ContentLength > MaxHTTPPayload {
		log.Println("bad content length:", r.ContentLength)
		return
	}
	parsedUrl, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		log.Println("bad url:", r.RequestURI, err)
		return
	}
	r.RequestURI = parsedUrl.Path

	isBlock := dnet.CheckBlocked(r)
	if isBlock {
		log.Println("is blocked by wx:", r.Host)
		w.Write([]byte(dnet.BlockHTML))
		return
	}
	proxy := GetProxyPort(r.Host)
	if proxy != nil {
		log.Println("using proxy:", r.Host)
		proxy.ServeHTTP(w, r)
		return
	}
	if !IsGoodSNI(r.Host) {
		log.Println("bad domain name:", r.Host)
		badDomainNameCounter.Inc()
		return
	}
	if strings.Index(r.RequestURI, "/memfile/") == 0 {
		MemFileHandler(w, r)
		return
	}
	if strings.Index(r.RequestURI, "/upload") == 0 {
		h.Uploader(w, r)
		return
	}
	if r.RequestURI == "/t/list" || r.RequestURI == "/t/list/" {
		h.MsgList(w, r)
		return
	}
	if r.RequestURI == "/t" {
		h.CreateMsg(w, r)
		return
	}
	if strings.Index(r.RequestURI, "/jump/") == 0 {
		Jump(w, r)
		return
	}
	if strings.Index(r.RequestURI, "/status/") == 0 {
		StatusCode(w, r)
		return
	}
	if strings.Index(r.RequestURI, "/t/list/") == 0 {
		h.MsgShow(w, r)
		return
	}
	if r.RequestURI == "/ip" {
		CFIPHandler(w, r)
		return
	}
    if r.RequestURI == "/vnstat" {
        VnstatHandler(w, r)
        return
    }
	if r.RequestURI == "/metrics" {
		promhttp.Handler().ServeHTTP(w, r)
		return
	}
	log.Println("serve file:", r.RequestURI)
	path := "/root/www"
	if r.Host == *FlagBlog {
		path = "/root/www/blogs"
	}
	if *FlagTestMode {
		path = "/Users/dilfish/go/src/github.com/dilfish/dynserver"
		if r.Host == *FlagBlog {
			path = "/Users/dilfish/go/src/github.com/dilfish/dynserver/blogs"
		}
	}
	d := http.Dir(path)
	f, err := d.Open(r.RequestURI)
	if err != nil {
		log.Println("open file error:", r.RequestURI, err)
	} else {
		fi, _ := f.Stat()
		fileSizeServedBytes.Set(float64(fi.Size()))
	}
	fs := http.FileServer(d)
	fs.ServeHTTP(w, r)
}

func IsGoodSNI(host string) bool {
    if *FlagBehindNginx {
        return true
    }
	h, _, err := net.SplitHostPort(host)
	if err == nil {
		host = h
	}
	for _, s := range SniList {
		if s == host {
			return true
		}
	}
	return false
}
