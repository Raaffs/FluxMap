package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
type CpmModel[T Analytic] struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (m *CpmModel[T])Insert(ctx context.Context, cpmValues []Cpm) error {
	query := `
	INSERT INTO Cpm (TaskID, EarliestStart, EarliestFinish, LatestStart, LatestFinish, SlackTime, CriticalPath, ParentProjectID, Dependencies)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (TaskID)
	DO UPDATE
	SET 
    EarliestStart = $2, 
    EarliestFinish = $3, 
    LatestStart = $4, 
    LatestFinish = $5, 
    SlackTime = $6, 
    CriticalPath = $7, 
    ParentProjectID = $8, 
    Dependencies = $9;
	`
	for _,cpm:=range cpmValues{
		_, err := m.DB.Exec(ctx, query, cpm.TaskID, cpm.EarliestStart, cpm.EarliestFinish, cpm.LatestStart, cpm.LatestFinish, cpm.SlackTime, cpm.CriticalPath,cpm.ParentProjectID,cpm.Dependencies)
		if err != nil {
			log.Printf("Error inserting data into Cpm table: %v", err)
			return err
		}
	}

	return nil
}

func(m *CpmModel[T])Exist()(bool,error){
	return true,nil
}

func(m *CpmModel[T])GetData(ctx context.Context,projectID int)([]*T,error){
	var cpmValues []*T
	query:=`SELECT TaskID, EarliestStart, EarliestFinish, LatestStart, LatestFinish, SlackTime, CriticalPath, ParentProjectID, Dependencies
	FROM cpm
	WHERE parentProjectID=$1

	`
	rows, err := m.DB.Query(ctx,query,projectID);if err!=nil{
		m.Errorlog.Printf("An error occurred while getting cpm values for projectID %v\n",err)
		if errors.Is(err,sql.ErrNoRows){
			return []*T{},ErrRecordNotFound
		}
		return []*T{},err
	}
	defer rows.Close()
	for rows.Next(){
		var cpm Cpm
		if err := rows.Scan(&cpm.TaskID,&cpm.EarliestStart,&cpm.EarliestFinish,&cpm.LatestStart,&cpm.LatestFinish,&cpm.SlackTime,&cpm.CriticalPath,&cpm.ParentProjectID, &cpm.Dependencies);err!=nil{
			m.Errorlog.Printf("An error occurred while scanning cpm values for projectID %v\n",err)
			return []*T{},err
		}
		var t *T
		var ok bool
		if t, ok = any(&cpm).(*T); !ok {
			return nil, fmt.Errorf("type assertion failed: cannot convert *cpm to %T", t)
		}
		cpmValues = append(cpmValues, t)
	}

	if rows.Err()!=nil{
		m.Errorlog.Printf("An error occurred while getting cpm values for projectID %v\n",err)
		return []*T{},err
	}

	return cpmValues,nil
}

func(m *CpmModel[T])InsertResult(ctx context.Context,projectID int, result Result)(error){
	query:=`
		INSERT into cpmResult(projectID, result)
		VALUES($1,$2)
		ON CONFLICT	(projectID)
		DO UPDATE SET
			result=EXCLUDED.result
	`
	_,err:=m.DB.Exec(ctx,query,projectID,result.Result);if err!=nil{
		return err
	}
	return nil
}

func(m *CpmModel[T])GetResult(ctx context.Context,projectID int)(Result,error){
	result:=Result{
		Result: map[string]any{},
	}
	query:=`SELECT result
	FROM cpmResult
	WHERE projectID=$1
	`
	if err :=m.DB.QueryRow(ctx,query,projectID).Scan(&result.Result);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return Result{},ErrRecordNotFound
		}
		m.Errorlog.Printf("An error occurred while getting cpm result for projectID %v\n",err)
		return Result{},err
	}
	return result,nil
}
