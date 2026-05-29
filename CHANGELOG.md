# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial scaffold — OOD compute adapter for single-node EC2 jobs, translating Open OnDemand job submissions to EC2 instance launches via a Launch Template.
- CLI commands: `submit` (launch an EC2 instance to run the job script), `status <instance-id>` (OOD-normalized status), `delete <instance-id>` (terminate), and `info <instance-id>` (full EC2 instance details as JSON).
- Unit and substrate integration tests.

### Changed
- Upgraded to substrate v0.45.2 and removed the `CreateLaunchTemplate` fallback.
