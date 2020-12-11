package domain

import "time"

// EcsService
type EcsService struct {
	Name         string
	CreatedAt    *time.Time
	LaunchType   string
	Status       string
	DesiredCount int64
	RunningCount int64
	PendingCount int64
}
