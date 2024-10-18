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
				log.Fatal("Handlers path is required")

				return
			}

			handlersPath := args[0]
			err := run.Deploy(handlersPath)
			if err != nil {
				log.Fatalf("Error during deployment: %v", err)
			}
		},
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy workerpool and redis containers",
		Run: func(_ *cobra.Command, _ []string) {
			err := run.Destroy()
			if err != nil {
				log.Fatalf("Error during destoy: %v", err)
			}
		},
	}

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push a task to the workerpool",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatal("Task type and payload required")
			}

			err := run.Push(args[0], args[1])
			if err != nil {
				log.Fatalf("Error during push: %v", err)
			}
		},
	}

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(pushCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
