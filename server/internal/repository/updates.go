package repository

import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)

type UpdateRepository interface {
    CreateUpdate(ctx context.Context, projectID int, msg string, createdBy string, targetType string, targetUsername *string) error
    GetProjectUpdates(ctx context.Context, projectID int, username string) ([]*models.Update, error)
    GetRecent(ctx context.Context, username string) ([]*models.Update, error) 	// Get the 7 most recent updates for a user
    GetAllUpdates(ctx context.Context, username string) ([]*models.Update, error)
    GetTotalUnreadUpdates(ctx context.Context, username string) (int, error)
    SetRead(ctx context.Context, username string) error     // Mark all updates as read for a user
}
