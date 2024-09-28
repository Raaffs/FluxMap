package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (t *TaskModel)Create(ctx context.Context, task Task)error{
	insert:=`INSERT INTO tasks(taskName,taskDescription,taskStatus,taskStartDate,taskDueDate,parentProjectID,assignedUsername) 
	VALUES($1,$2,$3,$4,$5,$6,$7)`
	_,err:=t.DB.Exec(ctx,insert,task.TaskName,task.TaskDescription,task.TaskStatus,task.TaskStartDate,task.TaskDueDate,task.ParentProjectID,task.AssignedUsername); if err!=nil{
		t.Errorlog.Println("Error creating project:",err)
		return err
	}
	return nil
}

func (t *TaskModel)UpdateManagerTask(ctx context.Context,task Task)error{
	query:=`
		UPDATE tasks
		SET taskname=$1, taskDescription=$2, taskStartDate=$3,taskDueDate=$4, AssignedUsername=$5, Approved=$6	
		WHERE taskID=$7
	`
	_,err:=t.DB.Exec(ctx,query,task.TaskName,task.TaskDescription,task.TaskStartDate,task.TaskDueDate,task.AssignedUsername,task.Approved,task.TaskID);if err!=nil{
		t.Errorlog.Println(err)
		return err
	}
	return nil
}

func(t *TaskModel)GetTasks(ctx context.Context,projectID int)([]*Task,error){
	var tasks []*Task
	query:=`
		SELECT taskID, taskName, taskDescription, taskStatus, taskStartDate, taskDueDate, parentProjectID, assignedUsername, Approved
		FROM tasks 
		WHERE parentProjectID=$1
	`
	rows,err:=t.DB.Query(ctx,query,projectID);if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*Task{},ErrRecordNotFound
		}
		return []*Task{},err
	}
	for rows.Next(){
		var task Task
		if err=rows.Scan(&task.TaskID,&task.TaskName,&task.TaskDescription,&task.TaskStatus,&task.TaskStartDate,&task.TaskDueDate,&task.ParentProjectID,&task.AssignedUsername,&task.Approved);err!=nil{
			return []*Task{},err
		}
		tasks=append(tasks, &task)
	}
	if err=rows.Err();err!=nil{
		return []*Task{},err
	}
	return tasks,err
}

func(t *TaskModel)GetTaskByID(ctx context.Context,taskID int)(Task,error){
	var task Task
	query:=`
		SELECT * FROM tasks
		WHERE taskID=$1
	`
	if err:=t.DB.QueryRow(ctx,query,taskID).Scan(&task.TaskID,&task.TaskName,&task.TaskDescription,&task.TaskStatus,&task.TaskStartDate,&task.TaskDueDate,&task.ParentProjectID,&task.AssignedUsername,&task.Approved);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return Task{},ErrRecordNotFound
		}
		return Task{},err
	}
	return task,nil
}
