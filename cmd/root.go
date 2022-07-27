/*
Copyright © 2022 Jace Walker <jc@jcwlkr.io>
*/
package cmd

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/tatsushid/go-fastping"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-netscanner",
	Short: "A very fast network scanner.",
	Long:  `A Gopher'd network scanner to return alive hosts within a given subnet.`,

	Run: func(cmd *cobra.Command, args []string) {
		argSubnet := cmd.Flag("subnet").Value.String()

		subnetAddresses := parseSubnet(argSubnet)

		var wg sync.WaitGroup
		for _, address := range subnetAddresses {
			// log.Println("Index:", index)
			wg.Add(1)
			go pingIP(address)

			defer wg.Done()
		}
		wg.Wait()
	},
}

func parseSubnet(subnet string) []net.IP {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(subnet)
	if err != nil {
		log.Fatal(err)
	}

	// convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)
	address := []net.IP{}

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		address = append(address, ip)
	}
	address = address[1 : len(address)-1]

	return address
}

func pingIP(ip net.IP) {
	p := fastping.NewPinger()
	// log.Println("Resolving IP Address", ip.String())
	ra, err := net.ResolveIPAddr("ip4:icmp", ip.String())
	if err != nil {
		log.Fatalln("Error:", err)
		os.Exit(1)
	}

	p.AddIPAddr(ra)
	// log.Println("Pinging", ra.String())
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("%s ALIVE, %v\n", addr.String(), rtt)
	}
	err = p.Run()
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-netscanner.yaml)")
	rootCmd.Flags().StringP("subnet", "s", "", "Subnet in CIDR format (eg 192.168.0.0/24)")
}