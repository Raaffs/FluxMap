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

func(u *UpdateModel)GetProjectUpdates(ctx context.Context, projectID int, username string) ([]*Update, error) {
	query := `
		SELECT 
		updates.id, 
		updates.projectid, 
		updates.msg, 
		updates.createdat, 
		updates.createdby, 
		updates.targettype, 
		updates.targetusername
		projects.projectName,
		projects.projectDescription

		FROM updates,projects

		WHERE updates.projectid = $1

		AND projects.projectID = updates.projectid
		AND (targettype = 'all' OR (targettype = 'user' AND targetusername = $2))

		ORDER BY createdat DESC
	`
	rows, err := u.DB.Query(ctx, query, projectID, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []*Update
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
			&upd.ProjectName,
			&upd.ProjectDescription,
		)
		if err != nil {
			return nil, err
		}
		updates = append(updates, &upd)
	}	
	if rows.Err()!=nil{
		return []*Update{},rows.Err()
	}
	return updates, nil
}

func (u *UpdateModel)GetAllUpdates(ctx context.Context, username string)([]*Update,error){
	var updates []*Update
	query := `
	SELECT 
	    updates.id, 
	    updates.projectid, 
	    updates.msg, 
	    updates.createdat, 
	    updates.createdby, 
	    updates.targettype, 
	    updates.targetusername,
	    projects.projectName, 
	    projects.projectDescription
	FROM 
	    updates
	JOIN 
	    projects ON updates.projectid = projects.projectid
	WHERE 
	    targettype='all' OR targetusername=$1
	ORDER BY 
	    createdat DESC;	`
		
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
			&upd.ProjectName,
			&upd.ProjectDescription,
		)
		if err!=nil{
			return []*Update{},err
		}
		updates=append(updates,&upd)
	}
	if rows.Err()!=nil{
		return []*Update{},rows.Err()
	}
	log.Println("updates in model: ",updates)
	return updates,nil
}

func (r *UpdateModel) CreateUpdate(ctx context.Context, projectID int, msg string, createdBy string, targetType string, targetUsername *string) ( error) {
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
