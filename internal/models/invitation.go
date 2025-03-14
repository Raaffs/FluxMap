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

func (i *InvitationModel)Invite(ctx context.Context,username string, projectID int)error{
	insert:=`
		INSERT INTO INVITATION(username,projectid)
		VALUES($1,$2)
	`
	_,err:=i.DB.Exec(ctx,insert,username,projectID);if err!=nil{
		return err
	}
	return nil
}

func (i *InvitationModel)GetInvitations(ctx context.Context, username string)([]*Invitation,error){
	var invitations []*Invitation
	get:=`SELECT * FROM INVITATION WHERE username=$1`
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

func (u *UserModel)ConfirmInvitation(ctx context.Context, projectID int,)error{
	query:=`
		UPDATE invitation
		SET accepted=true
		WHERE projectID=$1
	`
	_,err:=u.DB.Exec(ctx,query,projectID); if err!=nil{
		return err
	}
	return nil
}