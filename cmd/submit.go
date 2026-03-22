package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-ec2-adapter/internal/ec2"
	"github.com/scttfrdmn/ood-ec2-adapter/internal/ood"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Launch an EC2 instance to run an OOD job",
	Long:  "Reads a JSON job spec from stdin and launches an EC2 instance from the configured Launch Template.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec ood.JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		lt := launchTemplate
		if lt == "" {
			lt = spec.LaunchTemplate
		}
		if lt == "" {
			return fmt.Errorf("--launch-template is required (or set launch_template in job spec)")
		}

		// Encode the job script as user data
		userdata := base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\n" + spec.Script))

		ctx := context.Background()
		client, err := ec2.New(ctx, region)
		if err != nil {
			return err
		}

		instanceID, err := client.RunInstance(ctx, spec.JobName, lt, spec.InstanceType, userdata, nil)
		if err != nil {
			return err
		}

		fmt.Println(instanceID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
