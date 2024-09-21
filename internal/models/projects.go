package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
type ProjectModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (p *ProjectModel)Create(ctx context.Context, project Project)error{
	insert:=`INSERT INTO projects(projectName,projectDescription,projectDueDate,ownername)
	VALUES($1,$2,$3,$4)`
	_,err:=p.DB.Exec(ctx,insert,project.ProjectName,project.ProjectDescription,project.Ownername); if err!=nil{
		return err
	}
	return nil
}



func(p *ProjectModel)RetrieveAdmin(ctx context.Context,username string)([]*Project,error){
	var projects []*Project
	retrieve:=`SELECT projectName,projectDescription,projectDueDate FROM Projects WHERE ownername=$1`
	rows,err:=p.DB.Query(ctx,retrieve,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return projects,nil
		}
		return nil,err
	}
	defer rows.Close()

	for rows.Next(){
		var project Project  
		err=rows.Scan(&project.ProjectName,&project.ProjectDescription,&project.ProjectDueDate); if err!=nil{
			return nil,err
		}
		projects=append(projects,&project)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return projects,nil
}

func(p *ProjectModel)RetrieveManager(ctx context.Context,username string)([]*Project,error){
	var projects []*Project
	retrieve:=`SELECT projectName,projectDescription,projectDueDate FROM Projects WHERE ownername=$1`
	rows,err:=p.DB.Query(ctx,retrieve,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return projects,nil
		}
		return nil,err
	}
	for rows.Next(){
		var project Project
		err=rows.Scan(&project.ProjectName,&project.ProjectDescription,&project.ProjectDueDate); if err!=nil{
			return nil,err
		}
		projects=append(projects,&project)
	}
	return projects,nil
}


func(p *ProjectModel)RetrieveAssginedProjects(ctx context.Context,username string)([]*Project,error){
	var projects []*Project
	retrieve:=`SELECT projectName,projectDescription,projectDueDate FROM Projects WHERE ownername=$1`
	rows,err:=p.DB.Query(ctx,retrieve,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return projects,nil
		}
		return nil,err
	}
	for rows.Next(){
		var project Project
		err=rows.Scan(&project.ProjectName,&project.ProjectDescription,&project.ProjectDueDate); if err!=nil{
			return nil,err
		}
		projects=append(projects,&project)
	}
	return projects,nil
}
