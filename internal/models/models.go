package models

import (
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrRecordNotFound is returned when a movie record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
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
	Pert 		PertModel
	Cpm         CpmModel
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
		Pert: PertModel{
			DB:			db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
		Cpm: CpmModel{
			DB:			db,
			Infolog: 	infoLog,
			Errorlog: 	errorLog,
		},
	}
}