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
		INSERT INTO Cpm (TaskID, EarliestStart, EarliestFinish, LatestStart, LatestFinish, SlackTime, CriticalPath, ParentProjectID)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	for _,cpm:=range cpmValues{
		_, err := m.DB.Exec(ctx, query, cpm.TaskID, cpm.EarliestStart, cpm.EarliestFinish, cpm.LatestStart, cpm.LatestFinish, cpm.SlackTime, cpm.CriticalPath,cpm.ParentProjectID)
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

func(m *CpmModel[T])Get(ctx context.Context,projectID int)([]*T,error){
	var cpmValues []*T
	query:=`SELECT TaskID, EarliestStart, EarliestFinish, LatestStart, LatestFinish, SlackTime, CriticalPath, ParentProjectID 
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
		if err := rows.Scan(&cpm.TaskID,&cpm.EarliestStart,&cpm.EarliestFinish,&cpm.LatestStart,&cpm.LatestFinish,&cpm.SlackTime,&cpm.CriticalPath,&cpm.ParentProjectID);err!=nil{
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
