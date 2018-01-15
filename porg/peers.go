package main

import (
	"fmt"
	"math"
	"net"
)

const mask uint64 = 0xff
const bitsInByte = 8

func findPeers(net.IP, *net.IPNet) {

}

func main() {
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
