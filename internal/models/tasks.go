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
	insert:=`INSERT INTO tasks(taskName, taskDescription, taskStartDate, taskDueDate, parentProjectID, assignedUsername,taskCompletedDate,taskApprovedDate) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_,err:=t.DB.Exec(ctx,insert,task.TaskName,task.TaskDescription,task.TaskStartDate,task.TaskDueDate,task.ParentProjectID,task.AssignedUsername, task.TaskCompletedDate,task.TaskApprovedDate); if err!=nil{
		t.Errorlog.Println("Error creating project:",err)
		return err
	}
	return nil
}

func (t *TaskModel)UpdateTask(ctx context.Context,taskID int, status string)error{
	query:=`
	UPDATE tasks 
	SET 
	taskStatus=$1,
	taskCompletedDate=NOW()
	WHERE taskID=$2
	`
	_,err:=t.DB.Exec(ctx,query,status,taskID);if err!=nil{
		t.Errorlog.Println(err)
		return err
	}
	
	return nil
}


func (t *TaskModel)UpdateManagerTask(ctx context.Context,task Task)error{
	query:=`
	UPDATE tasks 
	SET taskname=$1, taskDescription=$2, taskStatus=$3, taskStartDate=$4, taskDueDate=$5, AssignedUsername=$6, Approved=$7, taskApprovedDate=$8
	WHERE taskID=$9
	`
	_,err:=t.DB.Exec(ctx,query,task.TaskName,task.TaskDescription,task.TaskStatus,task.TaskStartDate,task.TaskDueDate,task.AssignedUsername,task.Approved,task.TaskApprovedDate.Time.Format("2006-01-02"),task.TaskID);if err!=nil{
		t.Errorlog.Println(err)
		return err
	}
	return nil
}

func(t *TaskModel)GetTasks(ctx context.Context,projectID int)([]*Task,error){
	var tasks []*Task
	query:=`
		SELECT taskID, taskName, taskDescription, taskStatus, taskStartDate, taskDueDate, parentProjectID, assignedUsername, Approved, taskCompletedDate, taskApprovedDate
		FROM tasks 
		WHERE parentProjectID=$1
		ORDER BY taskID
	`
	rows,err:=t.DB.Query(ctx,query,projectID);if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*Task{},ErrRecordNotFound
		}
		return []*Task{},err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err=rows.Scan(&task.TaskID,&task.TaskName,&task.TaskDescription,&task.TaskStatus,&task.TaskStartDate,&task.TaskDueDate,&task.ParentProjectID,&task.AssignedUsername,&task.Approved, &task.TaskCompletedDate, &task.TaskApprovedDate);err!=nil{
			return []*Task{},err
		}
		tasks=append(tasks, &task)
	}
	if err=rows.Err();err!=nil{
		return []*Task{},err
	}
	return tasks,err
}

func(t *TaskModel)GetUserTasks(ctx context.Context,projectID int,username string)([]*Task,error){
	var tasks []*Task
	query:=`
		SELECT taskID, taskName, taskDescription, taskStatus, taskStartDate, taskDueDate, parentProjectID, assignedUsername, Approved, taskCompletedDate, taskApprovedDate
		FROM tasks 
		WHERE parentProjectID=$1
		AND assignedUsername=$2
		ORDER BY taskIDp
	`
	rows,err:=t.DB.Query(ctx,query,projectID,username);if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return []*Task{},ErrRecordNotFound
		}
		return []*Task{},err
	}
	defer rows.Close()
	for rows.Next(){
		var task Task
		if err=rows.Scan(&task.TaskID,&task.TaskName,&task.TaskDescription,&task.TaskStatus,&task.TaskStartDate,&task.TaskDueDate,&task.ParentProjectID,&task.AssignedUsername,&task.Approved, &task.TaskCompletedDate, &task.TaskApprovedDate);err!=nil{
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
