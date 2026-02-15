package repository

import (
	"context"

	"github.com/Raaffs/FluxMap/internal/models"
)
type InvitationRepository interface {
    Invite(ctx context.Context, username string, projectID int, role string) error
    GetPendingInvitations(ctx context.Context, username string) ([]*models.DisplayInvitations, error)
    ConfirmInvitation(ctx context.Context, status string, id int, username string) (int, error)
    Exist(ctx context.Context, username string, projectID int) (bool, error)
    AcceptedByUser(ctx context.Context, projectID int, username string) (bool, error)
    FetchConfirmedMembers(ctx context.Context, projectID int) ([]*string, error)
    GetInvitationByID(ctx context.Context, invitationID int) (*models.Invitation, error)
    GetTotalUnreadInvitation(ctx context.Context, username string) (int, error)
    SetRead(ctx context.Context, username string) error
    Delete(ctx context.Context, projectID int, username string) error
}
