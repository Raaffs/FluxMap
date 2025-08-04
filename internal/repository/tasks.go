package repository

import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)

type TaskRepository interface {
    Create(ctx context.Context, task models.Task) error
    UpdateStatus(ctx context.Context, projectID int, taskID int, status string) (string, error)
    ManagerAuthorizedUpdate(ctx context.Context, projectID int, task models.Task) error
    Get(ctx context.Context, projectID int) ([]*models.Task, error)
    GetAssignedToUser(ctx context.Context, projectID int, username string) ([]*models.Task, error)
    GetByID(ctx context.Context, taskID int) (models.Task, error)
    Archieve(ctx context.Context, projectID int, taskID int) error
    ReallocateUser(ctx context.Context, projectID int, removedUser, fallBackUser string) error
}
