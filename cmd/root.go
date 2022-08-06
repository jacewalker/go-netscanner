/*
Copyright © 2022 Jace Walker <jc@jcwlkr.io>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jacewalker/go-netscanner/cmd/ping"
	"github.com/jacewalker/go-netscanner/cmd/ports"
	"github.com/spf13/cobra"
)

var (
	subnet          string
	portsString     string
	commonPortCheck bool
	commonPorts     = []int{80, 443, 22}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-netscanner",
	Short: "A very fast network scanner.",
	Long:  `A Gopher'd network scanner to return alive hosts within a given subnet.`,

	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		subnetAddresses := ping.ParseSubnet(subnet)

		var wg sync.WaitGroup
		for _, address := range subnetAddresses {
			wg.Add(1)

			if portsString != "0" {
				go ports.ScanPorts(address, portsString)
			}

			go ping.PingIP(address, &wg)
		}
		wg.Wait()

		duration := time.Since(startTime).Truncate(1000000)
		fmt.Println("Duration:", duration)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&subnet, "subnet", "s", subnet, "Subnet in CIDR format (eg 192.168.0.0/24)")
	rootCmd.Flags().StringVarP(&portsString, "ports", "p", portsString, "Specify a single port or range (0-1000) of ports to scan")
	rootCmd.Flags().BoolVarP(&commonPortCheck, "commonports", "c", commonPortCheck, "Call -c to scan common ports (eg 80, 443, 22)")
}
