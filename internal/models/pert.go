package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PertModel[T  Analytic] struct {
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func(p *PertModel[T])Insert(ctx context.Context,PertValues []Pert)error{
	query:=`
		INSERT INTO Pert(ParentTaskID, PredecessorTaskID, Optimistic, Pessimistic, MostLikely, ParentProjectID)
		VALUES($1, $2, $3, $4, $5, $6)
	`
	for _,val :=range PertValues{
		_, err := p.DB.Exec(ctx,query,val.ParentTaskID, val.PredecessorTaskID,val.Optimistic,val.Pessimistic,val.MostLikely,val.ParentProjectID);if err!=nil{
			p.Errorlog.Printf("An error occurred wile inserting %v in pert table",val)
			return err
		}
	}
	return nil
}
func(p *PertModel[T])Exist()(bool,error){
	return true,nil
}
func(p *PertModel[T])Get(ctx context.Context,projectID int)([]*T,error){
	var pertValues []*T
	query:=`SELECT parentTaskID, predecessorTaskID, optimistic, pessimistic, mostLikely, parentProjectID 
	FROM pert
	WHERE parentProjectID=$1
	`
	rows, err := p.DB.Query(ctx,query,projectID);if err!=nil{
		p.Errorlog.Printf("An error occurred while getting pert values for projectID %v\n",err)
		if errors.Is(err,sql.ErrNoRows){
			return []*T{},ErrRecordNotFound
		}
		return []*T{},err
	}
	defer rows.Close()
	for rows.Next(){
		var pert Pert
		if err := rows.Scan(&pert.ParentTaskID,&pert.PredecessorTaskID,&pert.Optimistic,&pert.Pessimistic,&pert.MostLikely,&pert.ParentProjectID);err!=nil{
			p.Errorlog.Printf("An error occurred while scanning pert values for projectID %v\n",err)
			return []*T{},err
		}
		// Use a type assertion to convert pert to T
		var t *T
		var ok bool
		if t, ok = any(&pert).(*T); !ok {
			return nil, fmt.Errorf("type assertion failed: cannot convert *Pert to %T", t)
		}
		pertValues = append(pertValues, t)
	}

	if rows.Err()!=nil{
		p.Errorlog.Printf("An error occurred while getting pert values for projectID %v\n",err)
		return []*T{},err
	}

	return pertValues,nil
}

