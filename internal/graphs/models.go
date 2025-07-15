package graphs

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GraphModel struct {
	DB       *pgxpool.Pool
	Infolog  *log.Logger
	Errorlog *log.Logger
}

func NewGraphs(db *pgxpool.Pool) GraphModel {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return GraphModel{
		DB:       db,
		Infolog:  infoLog,
		Errorlog: errorLog,
	}
}