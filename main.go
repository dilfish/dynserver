package main

import (
	"flag"
	dnet "github.com/dilfish/tools/net"
	"log"
	"net/http"
	"strings"
	"time"
)

var FlagCert = flag.String("cert", "dev.ug.pem", "cert file path")
var FlagKey = flag.String("key", "dev.ug.key.pem", "priv key file path")
var FlagDomain = flag.String("domain", "ak.dev.ug", "domain name")
var FlagSNI = flag.String("sni", "ls.dev.ug,ls4.dev.ug", "sni list")

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	domain := *FlagDomain
	SniList = strings.Split(*FlagSNI, ",")
	if len(SniList) == 0 {
		panic("sni list is empty")
	}

	client := NewMongoClient("mongodb://localhost:27017", "msglist", "msg")
	if client == nil {
		panic("new mongo client error")
	}

	go Redirect()
	var h HttpsHandler
	h.C = client
	h.u = dnet.NewUploadService(
		"https://"+domain+"/ugc/",
		"/root/go/src/dynserver/ugc",
		"https://"+domain+"/upload",
		1024*1024*10,
		time.Hour*24*30, 5)
	err := http.ListenAndServeTLS(":443", *FlagCert, *FlagKey, &h)
	if err != nil {
		log.Println("listen and serve tls error:", err)
	}
}
