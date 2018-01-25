package main

import (
	"database/sql"
)

func InitMysql(config DbConfig)(*sql.DB,error) {
	db, err := sql.Open("mysql", config.Name+":"+config.Pwd+"@tcp("+config.Ip+")/"+config.DbName+"?charset=utf8")
	return db,err
}