package repository

import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)


type ProjectRepository interface {
	Create(ctx context.Context, project models.Project) error
	Exist(ctx context.Context, projectID int, username string) (bool, error)
	RetrieveProjectByID(ctx context.Context, id int) (models.Project, error)
	RetrieveAdminProjects(ctx context.Context, username string) ([]*models.Project, error)
	RetrieveManagerProjects(ctx context.Context, username string) ([]*models.Project, error)
	RetrieveAssignedProjects(ctx context.Context, username string) ([]*models.Project, error)
	AssignManager(ctx context.Context, manager string, projectID int) error
}
