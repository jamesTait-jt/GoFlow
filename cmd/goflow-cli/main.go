package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jamesTait-jt/goflow/cmd/goflow-cli/run"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "goflow",
		Short: "Goflow CLI tool to deploy workerpool and plugins using Docker",
	}

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy workerpool with Redis broker and compiled plugins",
		Run: func(cmd *cobra.Command, args []string) {
			err := run.Deploy()
			if err != nil {
				log.Fatalf("Error during deployment: %v", err)
			}
		},
	}

	// Add flags for the deploy command
	// deployCmd.Flags().StringVarP(&brokerType, "broker", "b", "redis", "Specify the broker type (default: redis)")

	// Add deploy command to root
	rootCmd.AddCommand(deployCmd)

	// Execute Cobra root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
