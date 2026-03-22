package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-ec2-adapter/internal/ec2"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <instance-id>",
	Short: "Terminate an EC2 job instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := ec2.New(ctx, region)
		if err != nil {
			return err
		}
		if err := client.TerminateInstance(ctx, args[0]); err != nil {
			return err
		}
		fmt.Printf("Instance %s terminated\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
