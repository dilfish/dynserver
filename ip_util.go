package main

import (
	dnet "github.com/dilfish/tools/net"
	"net"
)

func IPv6ToUint64(ipstr string) (uint64, uint64) {
	return dnet.IPv62Num(ipstr)
}

func IpToU32(ip net.IP) uint32 {
	return dnet.IP2Num(ip.String())
}

func U32ToStr(u uint32) string {
	return dnet.Num2IP(u)
}

func Uint64sToStr(prefix, postfix uint64) string {
	return dnet.Num2IPv6(prefix, postfix)
}
