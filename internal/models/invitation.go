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

func (i *InvitationModel)GetPendingInvitations(ctx context.Context, username string)([]*DisplayInvitations,error){
	var invitations []*DisplayInvitations
	get:=`
		select 
		projects.projectid,
	    projects.projectname, 
		projects.projectdescription, 
		projects.ownername, 
		invitation.id, 
		invitation.role, 
		invitation.username

		from projects, invitation 

		where projects.projectid=invitation.projectid

		and invitation.username=$1
		and invitation.status = 'pending'
	`
	rows,err:=i.DB.Query(ctx,get,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*DisplayInvitations{},ErrRecordNotFound
		}
		return []*DisplayInvitations{},err
	}

	defer rows.Close()

	for rows.Next(){
		var invitation DisplayInvitations
		err=rows.Scan(
			&invitation.ProjectID,
			&invitation.ProjectName,
			&invitation.ProjectDescription,
			&invitation.OwnerName,
			&invitation.InvitationID,
			&invitation.Role,
			&invitation.Username,
		);if err!=nil{
			return []*DisplayInvitations{},err
		}
		invitations=append(invitations, &invitation)
	}
	if err=rows.Err();err!=nil{
		return []*DisplayInvitations{},err
	}
	return invitations,nil
}

func (i *InvitationModel)ConfirmInvitation(ctx context.Context, status string, id int,username string)(int,error){
	query:=`
		UPDATE 	invitation
		SET 	status=$1
		WHERE 	id=$2
		AND		username=$3
		RETURNING id
	`

	var invitationID int
	err := i.DB.QueryRow(ctx, query, status, id,username).Scan(&invitationID)
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

func (i *InvitationModel) HasAcceptedInvitation(ctx context.Context, projectID int, username string) (bool, error) {
	query := `
		SELECT id
		FROM invitation
		WHERE projectid = $1
		  AND username = $2
		  AND status = 'accepted'
	`

	row := i.DB.QueryRow(ctx, query, projectID, username)
	var dummy int
	err := row.Scan(&dummy)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error ",err)
			return false, nil // not accepted
		}
		log.Println("false, error",err)
		return false, err // something went boom
	}

	return true, nil // yes, they accepted
}


func (i *InvitationModel)FetchConfirmedMembers(ctx context.Context, projectID int)([]*string,error){
	query:=`
		SELECT username FROM invitation
		WHERE projectID=$1
		AND status='accepted'
	`
	var users []*string
	rows,err:=i.DB.Query(ctx,query,projectID); if err!=nil{
		return []*string{},err
	}
	defer rows.Close()
	for rows.Next(){
		var user string
		err=rows.Scan(&user); if err!=nil{
			return []*string{},err
		}
		users=append(users, &user)
	}
	if rows.Err()!=nil{
		return []*string{},rows.Err()
	}
	return users,nil
}

func (i *InvitationModel)GetInvitationByID(ctx context.Context, invitationID int) (*Invitation, error) {
	query := `
		SELECT id, projectID, username, status, role
		FROM invitation
		WHERE id = $1
	`

	var invitation Invitation
	err := i.DB.QueryRow(ctx, query, invitationID).Scan(
		&invitation.ID,
		&invitation.ProjectID,
		&invitation.Username,
		&invitation.Status,
		&invitation.Role,
	)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (i *InvitationModel)GetTotalUnreadInvitation(ctx context.Context, username string)(int,error){
	var count int 

	query:=`
		SELECT COUNT(*)
		FROM invitation
		WHERE username=$1
		AND hasread=false
	`
	if err:=i.DB.QueryRow(ctx,query,username).Scan(&count);err!=nil{
		return -1,err
	}
	return count,nil
}

func (i *InvitationModel)SetRead(ctx context.Context, username string)(error){
	query:=`
		UPDATE invitation
		SET hasread=true
		WHERE username=$1 
	` 
	_,err:=i.DB.Exec(ctx,query,username); if err!=nil{
		return err
	}
	return nil
}