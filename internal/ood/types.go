// Package ood defines the OOD job spec types for the EC2 adapter.
package ood

// JobSpec is the OOD job submission payload.
type JobSpec struct {
	Script         string            `json:"script"`
	JobName        string            `json:"job_name"`
	LaunchTemplate string            `json:"launch_template,omitempty"`
	InstanceType   string            `json:"instance_type,omitempty"`
	Walltime       string            `json:"walltime,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	NativeSpecs    []string          `json:"native_specs,omitempty"`
}

// JobStatus maps EC2 instance states to OOD status strings.
type JobStatus struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	InstanceID string `json:"instance_id,omitempty"`
	ExitCode   int    `json:"exit_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
	StatusUnknown   = "undetermined"
)
