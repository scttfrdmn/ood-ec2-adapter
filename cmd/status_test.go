package cmd

import (
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestEc2StateToOod(t *testing.T) {
	tests := []struct {
		name  string
		state awstypes.InstanceStateName
		want  string
	}{
		{
			name:  "pending maps to queued",
			state: awstypes.InstanceStateNamePending,
			want:  "queued",
		},
		{
			name:  "running maps to running",
			state: awstypes.InstanceStateNameRunning,
			want:  "running",
		},
		{
			name:  "stopped maps to completed",
			state: awstypes.InstanceStateNameStopped,
			want:  "completed",
		},
		{
			name:  "terminated maps to completed",
			state: awstypes.InstanceStateNameTerminated,
			want:  "completed",
		},
		{
			name:  "shutting-down maps to completed",
			state: awstypes.InstanceStateNameShuttingDown,
			want:  "completed",
		},
		{
			name:  "unknown state maps to undetermined",
			state: awstypes.InstanceStateName("banana"),
			want:  "undetermined",
		},
		{
			name:  "empty string maps to undetermined",
			state: awstypes.InstanceStateName(""),
			want:  "undetermined",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ec2StateToOod(tc.state)
			if got != tc.want {
				t.Errorf("ec2StateToOod(%q) = %q, want %q", tc.state, got, tc.want)
			}
		})
	}
}
