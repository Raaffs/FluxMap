package models

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
)

type Model interface{
	*pgxpool.Pool
	*log.Logger
	*log.Logger
}

// Models struct is a single convenient container to hold and represent all our database models.
type Models struct {
	Users       UserModel
	Projects    ProjectModel
	Task      	TaskModel
	Pert 		PertModel[Pert]
	Cpm         CpmModel[Cpm]
}

func NewModels(db *pgxpool.Pool) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Users: UserModel{
			DB: 		db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
		Projects: ProjectModel{
			DB:			db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
		Task: TaskModel{
			DB:			db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
		Pert: PertModel[Pert]{
			DB:			db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
		Cpm: CpmModel[Cpm]{
			DB:			db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
	}
}