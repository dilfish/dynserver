package main

import (
	"errors"
	"log"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var ProxyDomainList []string
var ProxyPortList []string
var ProxyList []*httputil.ReverseProxy

func ParseProxy() error {
	if *FlagProxyDomain == "" && *FlagProxyPort == "" {
		return nil
	}
	domainList := strings.Split(*FlagProxyDomain, ",")
	portList := strings.Split(*FlagProxyPort, ",")
	if len(domainList) != len(portList) {
		log.Println("reverse port list does not match domain list", *FlagProxyDomain, *FlagProxyPort)
		return errors.New("reverse port does not match domain list")
	}
	log.Println("domain and port list are:", domainList, portList)
	for _, p := range portList {
		_, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			log.Println("bad port value:", p, err)
			return err
		}
		t, err := url.Parse("http://localhost:" + p)
		if err != nil {
			log.Println("bad proxy url:", p, err)
			return err
		}
		ProxyList = append(ProxyList, httputil.NewSingleHostReverseProxy(t))
		ProxyPortList = append(ProxyPortList, p)
	}
	ProxyDomainList = domainList
	return nil
}

func GetProxyPort(d string) *httputil.ReverseProxy {
	log.Println("get proxy:", d, ProxyDomainList)
    if *FlagBehindNginx {
        return nil
    }
	for idx, p := range ProxyDomainList {
		if d == p {
			log.Println("proxy to:", p, ProxyPortList[idx])
			return ProxyList[idx]
		}
	}
	return nil
}
