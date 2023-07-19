package main

import (
	"flag"
	dnet "github.com/dilfish/tools/net"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strings"
	"time"
)

var FlagCert = flag.String("cert", "dev.ug.pem", "cert file path")
var FlagKey = flag.String("key", "dev.ug.key.pem", "priv key file path")
var FlagDomain = flag.String("domain", "ak.dev.ug", "domain name")
var FlagSNI = flag.String("sni", "ls.dev.ug,ls4.dev.ug", "sni list")
var FlagProxyPort = flag.String("pp", "", "proxy port list")
var FlagProxyDomain = flag.String("pd", "", "proxy domain list")

const MaxHTTPPayload = 1024 * 1024 * 10

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	err := ParseProxy()
	if err != nil {
		log.Println("parse proxy error:", err)
		return
	}

	domain := *FlagDomain
	SniList = strings.Split(*FlagSNI, ",")
	if len(SniList) == 0 {
		panic("sni list is empty")
	}

	client := NewMongoClient("mongodb://localhost:27017", "msglist", "msg")
	if client == nil {
		panic("new mongo client error")
	}

	var h HttpsHandler
	h.C = client
	h.u = dnet.NewUploadService(
		"https://"+domain+"/ugc/",
		"/root/go/src/dynserver/ugc",
		"https://"+domain+"/upload",
		MaxHTTPPayload,
		time.Hour*24*30, 5)
	go Redirect(&h)
	prometheus.MustRegister(badDomainNameCounter)
	prometheus.MustRegister(requestDurationUs)
	prometheus.MustRegister(fileSizeServedBytes)
	err = http.ListenAndServeTLS(":443", *FlagCert, *FlagKey, &h)
	if err != nil {
		log.Println("listen and serve tls error:", err)
	}
}
