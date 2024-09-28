package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (t *TaskModel)Create(ctx context.Context, task Task)error{
	var u UserModel
	u.DB = t.DB

	exists,err:=u.Exist(ctx,task.AssignedUsername.String);if err!=nil{
		t.Errorlog.Println(err)
		return err
	}
	if !exists{
		return ErrRecordNotFound
	}

	insert:=`INSERT INTO tasks(taskName,taskDescription,taskStatus,taskStartDate,taskDueDate,parentProjectID,assignedUsername) 
	VALUES($1,$2,$3,$4,$5,$6,$7)`
	_,err=t.DB.Exec(ctx,insert,task.TaskName,task.TaskDescription,task.TaskStatus,task.TaskStartDate,task.TaskDueDate,task.ParentProjectID,task.AssignedUsername); if err!=nil{
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

