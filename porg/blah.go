package main

import (
	"log"
	"math"
	"net"
)

func main() {
	var hi, lo uint64
	var offset uint

	addr, network, err := net.ParseCIDR("fd00::/8")
	ip := make(net.IP, net.IPv6len)

	log.Printf("%v %v %v", addr, network, err)
	for lo = 0; lo < math.MaxUint64; lo++ {
		for hi = 0; hi < math.MaxUint64; lo++ {
			for offset = 0; offset < net.IPv6len/2; offset++ {
				ip[offset] = byte(((0xff << (offset * 8)) & lo) >> (offset * 8))
			}

			for offset = net.IPv6len / 2; offset < net.IPv6len; offset++ {
				ip[offset] = byte(((0xff << (offset * 8)) & hi) >> ((offset - 8) * 8))
			}

			log.Println(ip)
		}
	}

	log.Println(network.Contains(addr))
}
