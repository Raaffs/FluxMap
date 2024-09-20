package models

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
type ProjectModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}