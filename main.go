package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goflow",
		Short: "GoFlow CLI for task processing",
	}

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy GoFlow with specified broker and handlers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("deploy")
			fmt.Println(os.Getwd())
		},
	}

	rootCmd.AddCommand(deployCmd) //, showCmd, pushCmd, getCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func deploy(cmd *cobra.Command, args []string) {
// 	brokerType, _ := cmd.Flags().GetString("broker")
// 	handlersPath, _ := cmd.Flags().GetString("handlers")

// 	cmd := exec.Command("docker-c")
// }
