package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateModel struct{
	DB 			*pgxpool.Pool
	Infolog		*log.Logger
	Errorlog	*log.Logger
}

func(u *UpdateModel)GetProjectUpdates(ctx context.Context, projectID int, username string) ([]Update, error) {
	query := `
		SELECT id, projectid, msg, createdat, createdby, targettype, targetusername
		FROM updates
		WHERE projectid = $1
		AND (targettype = 'all' OR (targettype = 'user' AND targetusername = $2))
		ORDER BY createdat DESC
	`
	rows, err := u.DB.Query(ctx, query, projectID, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []Update
	for rows.Next() {
		var upd Update
		err := rows.Scan(
			&upd.ID,
			&upd.ProjectID,
			&upd.Msg,
			&upd.CreatedAt,
			&upd.CreatedBy,
			&upd.TargetType,
			&upd.TargetUsername,
		)
		if err != nil {
			return nil, err
		}
		updates = append(updates, upd)
	}
	return updates, rows.Err()
}

func (u *UpdateModel)GetAllUpdates(ctx context.Context, username string)([]*Update,error){
	var updates []*Update
	query := `
		SELECT * FROM Updates
		WHERE targettype='all' or targetusername=$1
		ORDER BY createdat DESC
	`
	
	rows,err:=u.DB.Query(
		ctx,
		query,
		username,
	)
	
	if err!=nil{
		return []*Update{},err
	}

	defer rows.Close()
	for rows.Next(){
		var upd Update
		err=rows.Scan(
			&upd.ID,
			&upd.ProjectID,
			&upd.Msg,
			&upd.CreatedAt,
			&upd.CreatedBy,
			&upd.TargetType,
			&upd.TargetUsername,
		)
		if err!=nil{
			return []*Update{},err
		}
		updates=append(updates,&upd)
	}
	if rows.Err()!=nil{
		return []*Update{},rows.Err()
	}
	return updates,nil
}

func (r *UpdateModel) CreateUpdate(ctx context.Context, projectID int, msg string, createdBy string, targetType string, targetUsername string) ( error) {
	query := `
		INSERT INTO updates (projectid, msg, createdat, createdby, targettype, targetusername)
		VALUES ($1, $2, NOW(), $3, $4, $5)
	`
	_,err := r.DB.Exec(ctx, query,
		projectID,
		msg,
		createdBy,
		targetType,
		targetUsername,
	)
	return err
}
