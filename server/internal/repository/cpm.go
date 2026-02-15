package repository
//this entire package is for future use
import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)

type CpmRepository[T models.Analytic] interface {
    Insert(ctx context.Context, cpmValues []models.Cpm) error
    Exist() (bool, error)
    GetData(ctx context.Context, projectID int) ([]*T, error)
    InsertResult(ctx context.Context, projectID int, result models.Result) error
    GetResult(ctx context.Context, projectID int) (models.Result, error)
}