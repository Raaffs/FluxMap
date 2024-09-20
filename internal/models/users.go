package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)
type UserModel struct{
	DB 			*pgxpool.Pool
	Infolog 	*log.Logger
	Errorlog	*log.Logger
}

func (u *UserModel)Insert(ctx context.Context, user User)error{
	fmt.Println("users here : ",user.Email,user.Password,user.Username)
	exist,err:=u.Exist(ctx,user.Username); if err!=nil{
		fmt.Println("err init",err)
		if !errors.Is(err,sql.ErrNoRows){
			return err
		}
	}
	if exist{
		return ErrAlreadyExist
	}
	insert:=`INSERT INTO users (username, email, hashedPassword)
        VALUES ($1, $2, $3)`
	_,err=u.DB.Exec(ctx,insert,user.Username,user.Email,user.HashedPassword);if err!=nil{
		u.Errorlog.Println(err)
		return err
	}
	return nil
}

func(u *UserModel)Exist(ctx context.Context,username string)(bool,error){
	selectQuery:=`SELECT username FROM users WHERE username=$1`
	var user string
	err:=u.DB.QueryRow(ctx,selectQuery,username).Scan(&user)
	if err!=nil{
		if err==sql.ErrNoRows{
			return false,nil
		}
		return false,err
	}
	return true,nil
}