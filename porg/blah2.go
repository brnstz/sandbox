package main

import (
	"log"
	"net"
)

func main() {
	ip := make(net.IP, net.IPv6len)
	ip[net.IPv6len-1] = 0
	ip[net.IPv6len-2] = 1
	log.Println(ip)
	log.Println(net.IPv6len)

}
