/*
Copyright © 2022 Jace Walker <jc@jcwlkr.io>
*/
package ping

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/tatsushid/go-fastping"
)

func ParseSubnet(subnet string) []net.IP {
	_, ipv4Net, err := net.ParseCIDR(subnet)
	if err != nil {
		log.Fatal("Error:", err)
	}

	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	finish := (start & mask) | (mask ^ 0xffffffff)
	address := []net.IP{}

	for i := start; i <= finish; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		address = append(address, ip)
	}
	address = address[1 : len(address)-1]

	return address
}

func PingIP(ip net.IP, wg *sync.WaitGroup, activeHosts *[]string) {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip.String())
	if err != nil {
		log.Fatalln("Error:", err)
		os.Exit(1)
	}

	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		// fmt.Printf("%s ALIVE, %v\n", addr.String(), rtt)
		*activeHosts = append(*activeHosts, addr.String())
	}
	err = p.Run()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer wg.Done()
}
