package main

import (
	"flag"
	"fmt"
	dnet "github.com/dilfish/tools/net"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
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
var FlagT = flag.String("td", "", "telegram file folder")
var FlagTB = flag.String("tb", "", "telegram base url")
var FlagNToken = flag.Bool("nt", false, "using new telegram token")
var FlagBlog = flag.String("b", "", "blog domain")
var FlagBehindNginx = flag.Bool("bn", false, "behind nginx")
var FlagBehindNginxPort = flag.Int("np", 10080, "behind nginx port")

const MaxHTTPPayload = 1024 * 1024 * 30
const TgToken = "1153923115:AAHUig2LQfApIF_Q-v5fn_fKgkCYhI15Flc"
const TgFSToken = "6676857975:AAHkSi5n0ywJPWXu8HDvet1_u5PJtDvRAnU"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	if *FlagT != "" {
		if *FlagNToken {
			go Telegram(*FlagT, TgFSToken)
		} else {
			go Telegram(*FlagT, TgToken)
		}
	}

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
	err = IPCheckInit("sorted.memip.txt", "sorted.memip.v6.txt")
	if err != nil {
		panic("init ip check error")
	}

	var h HttpsHandler
	h.C = client
	h.u = dnet.NewUploadService(
		slog.New(slog.NewTextHandler(os.Stdout)),
		"https://"+domain+"/ugc/",
		"/root/www/ugc",
		"https://"+domain+"/upload",
		MaxHTTPPayload,
		time.Hour*24*365*3, 5)

	if *FlagBehindNginx {
		log.Println("mode behind nginx, using:", *FlagBehindNginxPort)
		err = http.ListenAndServe(":"+strconv.FormatInt(int64(*FlagBehindNginxPort), 10), &h)
		if err != nil {
			log.Println("liste error:", err)
		}
		return
	}

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
