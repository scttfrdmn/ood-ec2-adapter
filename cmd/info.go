package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-ec2-adapter/internal/ec2"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <instance-id>",
	Short: "Print full EC2 instance details as JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := ec2.New(ctx, region)
		if err != nil {
			return err
		}
		instance, err := client.DescribeInstance(ctx, args[0])
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(instance)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
