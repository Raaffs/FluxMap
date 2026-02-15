package repository

import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user models.User) error
    Login(ctx context.Context, username, password string) error
    Exist(ctx context.Context, username string) (bool, error)
    IsManager(ctx context.Context, username string, projectID string) (bool, error)
    IsAdmin(ctx context.Context, username string, projectID string) (bool, error)
    Invite(ctx context.Context, username string, projectID int) error     // Invite a user to a project
}
