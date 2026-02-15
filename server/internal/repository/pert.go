package repository

import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)

type PertRepository[T models.Analytic] interface {
    // From ReadDatabase
    Exist() (bool, error)
    InsertResult(ctx context.Context, projectID int, result models.Result) error
    GetData(ctx context.Context, projectID int) ([]*T, error)
    GetResult(ctx context.Context, projectID int) (models.Result, error)

    // Additional Pert-specific methods
    Insert(ctx context.Context, pertValues []models.Pert) error
    UpdateResult(ctx context.Context, projectID int, result models.Result) error
}
