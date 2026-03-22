package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-ec2-adapter/internal/ec2"
	"github.com/scttfrdmn/ood-ec2-adapter/internal/ood"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <instance-id>",
	Short: "Get the OOD status of an EC2 job instance",
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

		js := ood.JobStatus{
			ID:         args[0],
			InstanceID: args[0],
			Status:     ec2StateToOod(instance.State.Name),
		}
		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func ec2StateToOod(s awstypes.InstanceStateName) string {
	switch s {
	case awstypes.InstanceStateNamePending:
		return ood.StatusQueued
	case awstypes.InstanceStateNameRunning:
		return ood.StatusRunning
	case awstypes.InstanceStateNameStopped, awstypes.InstanceStateNameTerminated:
		return ood.StatusCompleted
	case awstypes.InstanceStateNameShuttingDown:
		return ood.StatusCompleted
	default:
		return ood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
