package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)
type UserModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (u *UserModel) Create(ctx context.Context, user User) error {
    exist, err := u.Exist(ctx, user.Username)
    if err != nil {
        if !(errors.Is(err,sql.ErrNoRows)){
			return err
		}
    }

    if exist {
        return ErrAlreadyExist // Return if user exists
    }
    insert := `INSERT INTO users (username, email, hashedPassword) VALUES ($1, $2, $3)`
    _, err = u.DB.Exec(ctx, insert, user.Username, user.Email, user.HashedPassword)
    if err != nil {
        u.Errorlog.Println("error inserting user", err)
        return err
    }
    return nil
}

func(u *UserModel)Login(ctx context.Context, username, password string)(error){
	var hashedPassword string
	_,err:=u.Exist(ctx,username);if err!=nil{
		return err
	}
	query:=`SELECT hashedPassword from users WHERE username=$1`
	err=u.DB.QueryRow(ctx,query,username).Scan(&hashedPassword)
	if err!=nil{
		return err
	}
	if err:=bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password));err!=nil{
		return ErrInvalidCredential
	}
	
	return nil
}

func(u *UserModel)Exist(ctx context.Context,username string)(bool,error){
	selectQuery:=`SELECT username FROM users WHERE username=$1`
	var user string
	err:=u.DB.QueryRow(ctx,selectQuery,username).Scan(&user)
	if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return false,nil
		}
		return false,err
	}
	return true,nil
}

func(u *UserModel)IsManager(ctx context.Context,username string, projectID string)(bool,error){
	var isManager bool
    query := `SELECT COUNT(*) > 0 FROM managers WHERE manage = $1 AND projectID = $2`
    if err := u.DB.QueryRow(ctx,query, username, projectID).Scan(&isManager);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return false,nil
		}
		return false,err
	}
    return isManager,nil
}

func(u *UserModel)IsAdmin(ctx context.Context,username string, projectID string)(bool,error){
	var isAdmin bool
    query := `SELECT COUNT(*) > 0 FROM projects WHERE manager = $1 AND projectID = $2`
    if err := u.DB.QueryRow(ctx,query, username, projectID).Scan(&isAdmin);err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return false,nil
		}
		return false,err
	}
    return isAdmin,nil
}