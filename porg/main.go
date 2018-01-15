package main

import (
	"log"
	"net"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	MinPort     int      `default:"1000"`
	MaxPort     int      `default:"2000"`
	SubnetMasks []string `default:"10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fd00::/8"`
}

func main() {
	var (
		c   config
		err error
	)

	err = envconfig.Process("porg", &c)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.MinPort, c.MaxPort, c.SubnetMasks)

	for _, mask := range c.SubnetMasks {
		ip, ipnet, err := net.ParseCIDR(mask)
		if err != nil {
			log.Println("can't parse mask %v %v", mask, err)
			continue
		}
		allIPs := findPeers(ip, ipnet)
		for ip := range allIPs {
			log.Println(ip)
		}
	}
}
