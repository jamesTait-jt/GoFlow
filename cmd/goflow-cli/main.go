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
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("handlers path is required")

				return
			}

			handlersPath := args[0]
			err := run.Deploy(handlersPath)
			if err != nil {
				log.Fatalf("Error during deployment: %v", err)
			}
		},
	}

	// Add deploy command to root
	rootCmd.AddCommand(deployCmd)

	// Execute Cobra root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
