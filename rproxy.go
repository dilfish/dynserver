package main

import (
    "log"
    "errors"
    "net/http/httputil"
    "strings"
    "strconv"
    "net/url"
)

var ProxyDomainList []string
var ProxyPortList []string
var ProxyList []*httputil.ReverseProxy

func ParseProxy() error {
    domainList := strings.Split(*FlagProxyDomain, ",")
    portList := strings.Split(*FlagProxyPort, ",")
    if len(domainList) != len(portList) {
        log.Println("reverse port list does not match domain list", *FlagProxyDomain, *FlagProxyPort)
        return errors.New("reverse port does not match domain list")
    }
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
    for idx, p := range ProxyDomainList {
       if d == p {
            log.Println("proxy to:", p, ProxyPortList[idx])
            return ProxyList[idx]
       }
    }
    return nil
}
