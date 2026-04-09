/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Finds and lists all available printers on the device.",
	Long:  `The discover command scans the system to find all printers currently installed and connected. It lists both physical printers (connected via USB, LAN, or Wi-Fi) and virtual printers (like PDF writers), providing their official names as configured in the operating system. This is a crucial first step for any printing workflow, ensuring you can identify the correct printer before sending a print job.`,
	Run: func(cmd *cobra.Command, args []string) {
		Singleton.Tui.DiscoveryForm.Show()
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
