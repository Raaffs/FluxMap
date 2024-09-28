package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
type ProjectModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (p *ProjectModel)Create(ctx context.Context, project Project)error{
	insert:=`INSERT INTO projects(projectName,projectDescription,projectStartDate,projectDueDate,ownername)
	VALUES($1,$2,$3,$4,$5)`
	_,err:=p.DB.Exec(ctx,insert,project.ProjectName,project.ProjectDescription,project.ProjectStartDate,project.ProjectDueDate,project.Ownername); if err!=nil{
		p.Errorlog.Println("Error creating project:",err)
		return err
	}
	fmt.Println("here project")

	return nil
}

func(p *ProjectModel)RetrieveAdminProjects(ctx context.Context,username string)([]*Project,error){
	var projects []*Project
	retrieve:=`SELECT projectID,projectName,projectDescription,projectDueDate FROM Projects WHERE ownername=$1`
	rows,err:=p.DB.Query(ctx,retrieve,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			p.Errorlog.Println("err no rows: ",err)
			return []*Project{},nil
		}
		p.Errorlog.Println("Error retrieving projects: ",err)
		return []*Project{},err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project  
		err=rows.Scan(&project.ProjectID,&project.ProjectName,&project.ProjectDescription,&project.ProjectDueDate); if err!=nil{
			return []*Project{},err
		}
		projects=append(projects,&project)
	}
	fmt.Println("scanned all")
	if err = rows.Err(); err != nil {
		p.Errorlog.Println("Error reading rows: ",err)
		return []*Project{}, err
	}
	return projects,nil
}

func(p *ProjectModel)RetrieveManagerProjects(ctx context.Context,username string)([]*Project,error){
	var projects []*Project
	retrieve:=`SELECT projects.ProjectID,projects.projectName,projects.projectDescription,projects.projectDueDate 
	FROM Projects
	JOIN Managers ON projects.projectID=managers.projectID
	WHERE managers.managername=$1`

	rows,err:=p.DB.Query(ctx,retrieve,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*Project{},nil
		}
		return []*Project{},err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		err=rows.Scan(&project.ProjectID,&project.ProjectName,&project.ProjectDescription,&project.ProjectDueDate); if err!=nil{
			return []*Project{},err
		}
		projects=append(projects,&project)
	}
	return projects,nil
}


func(p *ProjectModel)RetrieveAssginedProjects(ctx context.Context,username string)([]*Project,error){
	var projects []*Project
	retrieve:=`SELECT projects.projectID,projects.projectName,projects.projectDescription,projects.projectDueDate 
	FROM projects 
	JOIN tasks ON projects.projectID=tasks.parentProjectID
	WHERE tasks.assignedUsername=$1`
	rows,err:=p.DB.Query(ctx,retrieve,username); if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*Project{},nil
		}
		return []*Project{},err
	}
	defer rows.Close()
	for rows.Next(){
		var project Project
		err=rows.Scan(&project.ProjectID,&project.ProjectName,&project.ProjectDescription,&project.ProjectDueDate); if err!=nil{
			return []*Project{},err
		}
		projects=append(projects,&project)
	}
	return projects,nil
}
func(p *ProjectModel)AssignManager(ctx context.Context,manager , projectID string)(error){
	query:=`INSERT INTO managers(managername,projectid)VALUES($1,$2)`	
	_,err:=p.DB.Exec(ctx,query,manager,projectID);if err!=nil{
		p.Errorlog.Println("Error assigning manager: ",err)
		return err
	}
	return nil
}


