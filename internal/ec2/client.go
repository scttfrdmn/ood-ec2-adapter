// Package ec2 wraps the AWS EC2 API for the OOD adapter.
package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Client wraps the AWS EC2 client.
type Client struct {
	svc    *ec2.Client
	region string
}

// New creates an EC2 client using the default AWS credential chain.
func New(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: ec2.NewFromConfig(cfg), region: region}, nil
}

// RunInstance launches an EC2 instance from a Launch Template and runs a job script via user data.
func (c *Client) RunInstance(ctx context.Context, jobName, launchTemplateID, instanceType, userdata string, tags map[string]string) (string, error) {
	input := &ec2.RunInstancesInput{
		MinCount: aws.Int32(1),
		MaxCount: aws.Int32(1),
		LaunchTemplate: &types.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(launchTemplateID),
		},
		UserData: aws.String(userdata),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         buildTags(jobName, tags),
			},
		},
	}
	if instanceType != "" {
		input.InstanceType = types.InstanceType(instanceType)
	}

	out, err := c.svc.RunInstances(ctx, input)
	if err != nil {
		return "", fmt.Errorf("ec2 RunInstances: %w", err)
	}
	if len(out.Instances) == 0 {
		return "", fmt.Errorf("no instances returned from RunInstances")
	}
	return aws.ToString(out.Instances[0].InstanceId), nil
}

// DescribeInstance returns the current state of an EC2 instance.
func (c *Client) DescribeInstance(ctx context.Context, instanceID string) (*types.Instance, error) {
	out, err := c.svc.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("ec2 DescribeInstances: %w", err)
	}
	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			if aws.ToString(i.InstanceId) == instanceID {
				return &i, nil
			}
		}
	}
	return nil, fmt.Errorf("instance %q not found", instanceID)
}

// TerminateInstance terminates an EC2 instance.
func (c *Client) TerminateInstance(ctx context.Context, instanceID string) error {
	_, err := c.svc.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("ec2 TerminateInstances: %w", err)
	}
	return nil
}

func buildTags(jobName string, extra map[string]string) []types.Tag {
	tags := []types.Tag{
		{Key: aws.String("Name"), Value: aws.String(jobName)},
		{Key: aws.String("OodJob"), Value: aws.String(jobName)},
	}
	for k, v := range extra {
		tags = append(tags, types.Tag{Key: aws.String(k), Value: aws.String(v)})
	}
	return tags
}
