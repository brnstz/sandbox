package main

import (
	"log"
	"math"
	"net"
)

const mask uint64 = 0xff
const bitsInByte = 8

func main() {
	var hi, lo uint64
	var offset int
	var shiftBits uint

	addr, network, err := net.ParseCIDR("fd00::/8")
	ip := make(net.IP, net.IPv6len)

	log.Printf("%v %v %v", addr, network, err)

	// Use two int64s to represent an ipv6 address
	for hi = 0; hi < math.MaxUint64; hi++ {
		for lo = 0; lo < math.MaxUint64; lo++ {

			// For each combination of lo and hi, convert this to an array
			// of bytes

			// Starting with ip[16 - 1], use lo to set
			// the latter bytes of the ip address. For example:
			// If lo=1, then offset[15] = 1 and the rest of the bytes
			// are 0.
			// If lo=256, then offset[15] = 0 and offset[14] = 1
			for offset = net.IPv6len - 1; offset >= net.IPv6len/2; offset-- {
				// Set the value of the byte as this offset
				shiftBits = ((net.IPv6len - 1) - uint(offset)) * bitsInByte
				ip[offset] = byte(((mask << shiftBits) & lo) >> shiftBits)
				log.Printf("lo offset: %v shiftbits: %v value: %v", offset, shiftBits, ip[offset])
				/*
					byteValue = byte(((mask << shiftBits & lo)) >> shiftBits
					ip[offset] = byteValue
				*/
			}

			for offset = (net.IPv6len / 2) - 1; offset >= 0; offset-- {
				// Set the value of the byte as this offset
				shiftBits = ((net.IPv6len - 1) - (net.IPv6len / 2) - uint(offset)) * bitsInByte
				log.Printf("hi offset: %v shiftbits: %v", offset, shiftBits)
				/*
					byteValue = byte(((mask << shiftBits & lo)) >> shiftBits
					ip[offset] = byteValue
				*/
			}

			/*
				for offset = net.IPv6len / 2; offset < net.IPv6len; offset++ {
					ip[offset] = byte(((0xff << ((offset - 8) * 8)) & hi) >> ((offset - 8) * 8))
				}

				for offset = 0; offset < net.IPv6len/2; offset++ {
					ip[offset] = byte(((0xff << (offset * 8)) & lo) >> (offset * 8))
				}
			*/

			//log.Println(ip)
		}
	}

	//ip[net.IPv6len-1] = 1
	//log.Println(ip)

	//log.Println(network.Contains(addr))
}
