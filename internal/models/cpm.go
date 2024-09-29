package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
type CpmModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (m *CpmModel)Insert(ctx context.Context, cpmValues []Cpm) error {
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

func(m *CpmModel)Get(ctx context.Context,projectID int)([]*Cpm,error){
	var cpmValues []*Cpm
	query:=`SELECT TaskID, EarliestStart, EarliestFinish, LatestStart, LatestFinish, SlackTime, CriticalPath, ParentProjectID 
	FROM cpm
	WHERE parentProjectID=$1
	`
	rows, err := m.DB.Query(ctx,query,projectID);if err!=nil{
		m.Errorlog.Printf("An error occurred while getting pert values for projectID %v\n",err)
		if errors.Is(err,sql.ErrNoRows){
			return []*Cpm{},ErrRecordNotFound
		}
		return []*Cpm{},err
	}
	defer rows.Close()
	for rows.Next(){
		var cpm Cpm
		if err := rows.Scan(&cpm.TaskID,&cpm.EarliestStart,&cpm.EarliestFinish,&cpm.LatestStart,&cpm.LatestFinish,&cpm.SlackTime,&cpm.CriticalPath,&cpm.ParentProjectID);err!=nil{
			m.Errorlog.Printf("An error occurred while scanning pert values for projectID %v\n",err)
			return []*Cpm{},err
		}
		cpmValues=append(cpmValues, &cpm)
	}

	if rows.Err()!=nil{
		m.Errorlog.Printf("An error occurred while getting pert values for projectID %v\n",err)
		return []*Cpm{},err
	}

	return cpmValues,nil
}
