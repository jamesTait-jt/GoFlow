package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
			dockerfilePath := filepath.Join("dockerfiles", "Dockerfile.workerpool")
			if err := checkDockerfile(dockerfilePath); err != nil {
				fmt.Println("Error building workerpool image:", err)
			}

			if err := buildWorkerpoolImage(dockerfilePath); err != nil {
				fmt.Println("Error building workerpool image:", err)
				return
			}
		},
	}

	rootCmd.AddCommand(deployCmd) //, showCmd, pushCmd, getCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildWorkerpoolImage(dockerfilePath string) error {
	cmd := exec.Command("docker", "build", "-t", "goflow-workerpool", "-f", dockerfilePath, ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build workerpool image: %w", err)
	}
	return nil
}

// checkDockerfile checks if the Dockerfile exists at the given path.
func checkDockerfile(dockerfilePath string) error {
	// Use os.Stat to check if the file exists
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		// Return a descriptive error if the Dockerfile does not exist
		return fmt.Errorf("Dockerfile not found at: %s", dockerfilePath)
	}
	// Return nil if the file exists
	return nil
}

// func deploy(cmd *cobra.Command, args []string) {
// 	brokerType, _ := cmd.Flags().GetString("broker")
// 	handlersPath, _ := cmd.Flags().GetString("handlers")

// 	cmd := exec.Command("docker-c")
// }
