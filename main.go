package main

import (
	"flag"
	"fmt"
	dnet "github.com/dilfish/tools/net"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

var FlagCert = flag.String("cert", "dev.ug.pem", "cert file path")
var FlagKey = flag.String("key", "dev.ug.key.pem", "priv key file path")
var FlagDomain = flag.String("domain", "ak.dev.ug", "domain name")
var FlagSNI = flag.String("sni", "ls.dev.ug,ls4.dev.ug", "sni list")
var FlagProxyPort = flag.String("pp", "", "proxy port list")
var FlagProxyDomain = flag.String("pd", "", "proxy domain list")
var FlagTestMode = flag.Bool("t", false, "test mode")
var FlagV = flag.Bool("v", false, "print version info")

const MaxHTTPPayload = 1024 * 1024 * 30

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	if *FlagV {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Println("no build info")
			return
		}
		for _, setting := range bi.Settings {
			fmt.Println(setting.Key, setting.Value)
		}
		return
	}

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
	addr := ":443"
	if *FlagTestMode {
		addr = ":11443"
	}
	err = http.ListenAndServeTLS(addr, *FlagCert, *FlagKey, &h)
	if err != nil {
		log.Println("listen and serve tls error:", err)
	}
}
