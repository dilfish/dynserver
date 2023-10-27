package main

import (
	"log"
	"net"
	"strings"
)

var IP6C *IPv6Manager
var IP4C *RangerMap

func IPCheckInit(v4, v6 string) error {
	IP4C = NewRangerMap()
	err := IP4C.ReadFile(v4)
	if err != nil {
		log.Println("read file error:", err)
		return err
	}
	IP4C.Manage()

	ip6m := &IPv6Manager{}
	ip6m.ReadV6File(v6)
	ip6m.SortV6()
	IP6C = ip6m
	return nil
}

func IsGoodIP(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		log.Println("bad ip:", ipstr)
		return false
	}
	var view string
	// ipv6
	if ip.To4() == nil {
		view = IP6C.FindV6(ipstr)
	} else {
		view = IP4C.Find(IpToU32(ip))
	}
	array := strings.Split(view, "-")
	if len(array) != 6 {
		log.Println("bad ip:", ipstr, view)
		return false
	}
	log.Println("view info:", ipstr, view)
	// 城市-省份-大区-ISP-国家-大洲
	if array[4] != "中国" {
		return false
	}
	return true
}
