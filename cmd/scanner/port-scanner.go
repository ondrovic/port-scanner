package main

import (
	"log"
	"os"
	"port-scanner/internal/models"
	"port-scanner/internal/scanner"
	"port-scanner/internal/utils"

	"github.com/spf13/cobra"
)

func main() {
	var flags models.CommandFlags
	var rootCmd = &cobra.Command{
		Use:   "port-scanner",
		Short: "Cli Port Scanner",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.Host == "" {
				cmd.Help()
				os.Exit(1)
			}
			result := scanner.ScanPorts(flags.Host, flags.StartPort, flags.EndPort)
			utils.PrintResultAsJSON(result)
		},
	}

	rootCmd.Flags().StringVarP(&flags.Host, "target", "t", "", "Target you wish to scan")
	rootCmd.Flags().IntVarP(&flags.StartPort, "start-port", "s", 1, "Starting port number to scan")
	rootCmd.Flags().IntVarP(&flags.EndPort, "end-port", "e", 65535, "Ending port number to scan")
	
	rootCmd.MarkFlagsRequiredTogether("target")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
