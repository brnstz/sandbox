package main

import (
	"fmt"
	"math"
	"net"
)

const mask uint64 = 0xff
const bitsInByte = 8

// FIXME: test code
func printPeers() {
	var hi, lo uint64
	var offset int
	var shiftBits uint

	ip := make(net.IP, net.IPv6len)
	i := 0

	// Use two int64s to represent an ipv6 address
	for hi = 0; hi < math.MaxUint64; hi++ {
		for lo = 0; lo < math.MaxUint64; lo++ {

			// For each combination of lo and hi, convert this to an array
			// of bytes representing an IPv6 address (net.IP).

			// Starting with ip[16 - 1], use lo to set
			// the latter bytes of the ip address. For example:
			// If lo=1, then offset[15] = 1 and the rest of the bytes
			// are 0.
			// If lo=256, then offset[15] = 0 and offset[14] = 1
			for offset = net.IPv6len - 1; offset >= net.IPv6len/2; offset-- {
				shiftBits = ((net.IPv6len - 1) - uint(offset)) * bitsInByte
				ip[offset] = byte(((mask << shiftBits) & lo) >> shiftBits)
			}

			// Do the same but for the hi bits
			for offset = (net.IPv6len / 2) - 1; offset >= 0; offset-- {
				shiftBits = ((net.IPv6len - 1) - (net.IPv6len / 2) - uint(offset)) * bitsInByte
				ip[offset] = byte(((mask << shiftBits) & hi) >> shiftBits)
			}

			if i%*n == 0 {
				fmt.Println(ip)
			}
			i++
		}
	}
}

/*
FIXME: more test code
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
*/
