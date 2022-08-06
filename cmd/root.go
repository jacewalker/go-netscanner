/*
Copyright © 2022 Jace Walker <jc@jcwlkr.io>
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jacewalker/go-netscanner/cmd/ping"
	"github.com/jacewalker/go-netscanner/cmd/ports"
	"github.com/spf13/cobra"
)

var (
	subnetFlag     string
	portsFlag      string
	commonPortFlag bool
	commonPorts    = []int{80, 443, 22}
)

type Target struct {
	Hosts       []net.IP
	Ports       []int
	ActivePorts []int
	ActiveHosts []string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-netscanner",
	Short: "A very fast network scanner.",
	Long:  `A Gopher'd network scanner to return alive hosts and open ports within a given subnet.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}

		startTime := time.Now()

		target := Target{}
		target.Hosts = ping.ParseSubnet(subnetFlag)
		if portsFlag != "" {
			target.Ports = ports.ConvertPortsStringToSlice(portsFlag)
		}

		var wg sync.WaitGroup

		fmt.Println("\n################\nOpen Ports:")
		for _, address := range target.Hosts {
			wg.Add(1)

			if portsFlag != "0" {
				go ports.ScanPorts(address, target.Ports)
			} else if commonPortFlag {
				go ports.ScanPorts(address, commonPorts)
			}

			go ping.PingIP(address, &wg, &target.ActiveHosts)
		}
		wg.Wait()

		fmt.Println("\n################\nActive Hosts:")
		for _, host := range target.ActiveHosts {
			fmt.Println(host)
		}

		duration := time.Since(startTime).Truncate(1000000)
		fmt.Println("Duration:", duration)

		fmt.Println("\n################\nActive Ports:", target.ActivePorts)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&subnetFlag, "subnet", "s", subnetFlag, "(Required) Subnet in CIDR format (eg 192.168.0.0/24)")
	rootCmd.Flags().StringVarP(&portsFlag, "ports", "p", portsFlag, "Specify a single port or range (0-1000) of ports to scan")
	rootCmd.Flags().BoolVarP(&commonPortFlag, "commonports", "c", commonPortFlag, "Call -c to scan common ports (eg 80, 443, 22)")
}
