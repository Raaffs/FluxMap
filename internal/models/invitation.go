package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InvitationModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog 	*log.Logger
}

func (i *InvitationModel)Invite(ctx context.Context,username string, projectID int,role string)error{
	insert:=`
		INSERT INTO INVITATION(username,projectid,role)
		VALUES($1,$2,$3)
	`
	_,err:=i.DB.Exec(ctx,insert,username,projectID,role);if err!=nil{
		return err
	}
	return nil
}

func (i *InvitationModel)GetInvitations(ctx context.Context, username string)([]*Invitation,error){
	var invitations []*Invitation
	get:=`
		select projects.projectid, projects.projectname, projects.projectdescription, projects.ownername, invitation.id, invitation.role, invitation.username
		from projects, invitation 
		where projects.projectid=invitation.projectid
		and invitation.username=$1
	`
	rows,err:=i.DB.Query(ctx,get,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*Invitation{},ErrRecordNotFound
		}
		return []*Invitation{},err
	}

	defer rows.Close()

	for rows.Next(){
		var invitation Invitation
		err=rows.Scan(&invitation.ID,&invitation.Username,&invitation.ProjectID,&invitation.Accepted);if err!=nil{
			return []*Invitation{},err
		}
		invitations=append(invitations, &invitation)
	}
	if err=rows.Err();err!=nil{
		return []*Invitation{},err
	}
	return invitations,nil
}

func (i *InvitationModel)ConfirmInvitation(ctx context.Context, projectID int,username string)(int,error){
	query:=`
		UPDATE 	invitation
		SET 	accepted=true
		WHERE 	projectID=$1
		AND		username=$2
		RETURNING id
	`
	var invitationID int
	err := i.DB.QueryRow(ctx, query, projectID).Scan(&invitationID)
	if err != nil {
		return 0, err
	}
	return invitationID,nil
}

func (i *InvitationModel) Exist(ctx context.Context, username string, projectID int) (bool, error) {
	selectQuery := `SELECT EXISTS(SELECT 1 FROM invitation WHERE username = $1 AND projectID=$2)`
	var exists bool
	err := i.DB.QueryRow(ctx, selectQuery, username, projectID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (i *InvitationModel) GetInvitationByID(ctx context.Context, invitationID int) (*Invitation, error) {
	query := `
		SELECT id, projectID, username, accepted, role
		FROM invitation
		WHERE id = $1
	`

	var invitation Invitation
	err := i.DB.QueryRow(ctx, query, invitationID).Scan(
		&invitation.ID,
		&invitation.ProjectID,
		&invitation.Username,
		&invitation.Accepted,
		&invitation.Role,
	)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}


//saving this query for later
// select projects.projectid, projects.projectname, projects.projectdescription, projects.ownername, invitation.id, invitation.role, invitation.username
// from projects, invitation 
// where projects.projectid=invitation.projectid