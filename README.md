# ood-ec2-adapter

[![CI](https://github.com/scttfrdmn/ood-ec2-adapter/actions/workflows/ci.yml/badge.svg)](https://github.com/scttfrdmn/ood-ec2-adapter/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/ood-ec2-adapter)](https://goreportcard.com/report/github.com/scttfrdmn/ood-ec2-adapter)
[![codecov](https://codecov.io/gh/scttfrdmn/ood-ec2-adapter/branch/main/graph/badge.svg)](https://codecov.io/gh/scttfrdmn/ood-ec2-adapter)
[![Go Reference](https://pkg.go.dev/badge/github.com/scttfrdmn/ood-ec2-adapter.svg)](https://pkg.go.dev/github.com/scttfrdmn/ood-ec2-adapter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

OOD compute adapter for single-node EC2 jobs. Translates Open OnDemand job submissions to EC2 instance launches using a Launch Template.

## Commands

| Command | Description |
|---------|-------------|
| `submit` | Launch an EC2 instance to run the job script |
| `status <instance-id>` | Get OOD-normalized status of the EC2 instance |
| `delete <instance-id>` | Terminate the EC2 instance |
| `info <instance-id>` | Print full EC2 instance details as JSON |

## Usage

```bash
# Submit a job
echo '{"job_name":"myjob","script":"#!/bin/bash\necho hello > /tmp/out.txt"}' | \
  ood-ec2-adapter submit \
    --launch-template lt-xxxxxxxxxxxxxxx \
    --region us-east-1

# Check status
ood-ec2-adapter status i-xxxxxxxxxxxxxxxxx

# Terminate
ood-ec2-adapter delete i-xxxxxxxxxxxxxxxxx
```

## OOD Cluster Config

```yaml
# /etc/ood/config/clusters.d/aws-ec2.yml
---
v2:
  metadata:
    title: "AWS EC2 Compute"
  job:
    adapter: "adapter_script"
    submit_host: "localhost"
    submit:
      script: "/usr/local/lib/ood-adapters/ood-ec2-adapter"
      args:
        - submit
        - "--launch-template=lt-xxxxxxxxxxxxxxx"
        - "--region=us-east-1"
```

## Infrastructure

Terraform in `aws-openondemand` with `adapters_enabled = ["ec2"]` provisions:
- IAM policy on the OOD instance role for `ec2:RunInstances`, `ec2:TerminateInstances`, etc.
- A Launch Template should be created separately and referenced here.
