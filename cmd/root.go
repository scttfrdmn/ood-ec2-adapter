package cmd

import (
	"github.com/spf13/cobra"
)

var (
	region         string
	launchTemplate string
)

var rootCmd = &cobra.Command{
	Use:   "ood-ec2-adapter",
	Short: "OOD compute adapter for EC2 single-node jobs",
	Long:  "Translates Open OnDemand job submissions to EC2 instance launches.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&launchTemplate, "launch-template", "", "EC2 Launch Template ID")
}
