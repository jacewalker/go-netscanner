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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-netscanner",
	Short: "A very fast network scanner.",
	Long:  `A Gopher'd network scanner to return alive hosts within a given subnet.`,

	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		subnet := cmd.Flag("subnet").Value.String()
		portsString := cmd.Flag("ports").Value.String()

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
	rootCmd.Flags().StringP("subnet", "s", "", "Subnet in CIDR format (eg 192.168.0.0/24)")
	rootCmd.Flags().StringP("ports", "p", "0", "Specify a single port or range (0-1000) of ports to scan")
}
