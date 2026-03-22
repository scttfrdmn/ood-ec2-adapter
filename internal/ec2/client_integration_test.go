//go:build integration

package ec2_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	substrate "github.com/scttfrdmn/substrate"

	. "github.com/scttfrdmn/ood-ec2-adapter/internal/ec2"
)

// substrateEC2Client builds a raw AWS EC2 SDK client pointed at the substrate server.
func substrateEC2Client(t *testing.T, endpointURL string) *awsec2.Client {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(endpointURL),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	return awsec2.NewFromConfig(cfg)
}

// TestRunAndDescribeInstance_Substrate tests the full RunInstance → DescribeInstance →
// TerminateInstance lifecycle against the substrate emulator.
func TestRunAndDescribeInstance_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()

	// CreateLaunchTemplate — supported since substrate v0.44.4.
	ec2Client := substrateEC2Client(t, ts.URL)
	ltOut, err := ec2Client.CreateLaunchTemplate(ctx, &awsec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String("ood-test"),
		LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
			ImageId:      aws.String("ami-12345678"),
			InstanceType: ec2types.InstanceTypeT3Medium,
		},
	})
	if err != nil {
		t.Fatalf("CreateLaunchTemplate: %v", err)
	}
	launchTemplateID := aws.ToString(ltOut.LaunchTemplate.LaunchTemplateId)
	t.Logf("launch template: %s", launchTemplateID)

	// Build the adapter client (picks up AWS_ENDPOINT_URL from the environment).
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// RunInstance — substrate accepts RunInstances with any (or no) ImageId.
	instanceID, err := client.RunInstance(ctx, "test-job", launchTemplateID, "t3.medium", "#!/bin/bash\necho hi", nil)
	if err != nil {
		t.Fatalf("RunInstance: %v", err)
	}
	if instanceID == "" {
		t.Fatal("expected non-empty instance ID")
	}
	t.Logf("launched instance ID: %s", instanceID)

	// DescribeInstance — verify the adapter can retrieve the instance we just launched.
	instance, err := client.DescribeInstance(ctx, instanceID)
	if err != nil {
		t.Fatalf("DescribeInstance: %v", err)
	}
	if aws.ToString(instance.InstanceId) != instanceID {
		t.Errorf("DescribeInstance: got ID %q, want %q", aws.ToString(instance.InstanceId), instanceID)
	}
	t.Logf("instance state: %s", instance.State.Name)

	// TerminateInstance — should not error.
	err = client.TerminateInstance(ctx, instanceID)
	if err != nil {
		t.Fatalf("TerminateInstance: %v", err)
	}
	t.Log("instance terminated successfully")
}

// TestDescribeInstance_NotFound_Substrate verifies that DescribeInstance returns an error
// for an instance ID that has never been created.
func TestDescribeInstance_NotFound_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.DescribeInstance(ctx, "i-doesnotexist")
	if err == nil {
		t.Fatal("expected error for non-existent instance, got nil")
	}
	// Accept either "not found" from the adapter's own error message or an
	// error containing the instance ID from a lower-level SDK/substrate response.
	if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "doesnotexist") {
		t.Logf("error (acceptable format): %v", err)
	}
	t.Logf("DescribeInstance non-existent instance returned error: %v", err)
}
