package tasks

import (
	"context"

	"github.com/nanoteck137/refinery/service"
)

var _ service.Task = (*AuthCleanupTask)(nil)

type AuthCleanupTask struct {
	authService *service.AuthService
}

func NewAuthCleanupTask(authService *service.AuthService) *AuthCleanupTask {
	return &AuthCleanupTask{
		authService: authService,
	}
}

func (j *AuthCleanupTask) Name() string {
	return AuthCleanup
}

func (j *AuthCleanupTask) Schedule() string {
	return "@every 30m"
}

func (j *AuthCleanupTask) Run(ctx context.Context) error {
	j.authService.Cleanup()
	return nil
}
