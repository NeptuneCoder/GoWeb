package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"github.com/yanghai23/GoLib/atdb"
	"github.com/yanghai23/GoLib/aterr"
)

func main() {
	config,err:= atdb.ReadDbConfig("/Users/yh/AndroidStudioProjects/GoWeb/config/logDBConfig.json")
	if err != nil {
		fmt.Println("..............ttttt")
		return
	}

	db, err := atdb.InitMysql(*config)
	aterr.CheckErr(err)

	insertData(db)

	//
	queryData(db)
	//////删除数据
	//stmt, err := db.Prepare("delete from userinfo where uid=?")
	//checkErr(err)
	//res, err := stmt.Exec("12")
	//checkErr(err)

	//affect, err := res.RowsAffected()
	//checkErr(err)
	defer func() {
		fmt.Println("..............")
		db.Close()
	}()
	fmt.Println("test.....")

}
func insertData(db *sql.DB) {
		//插入数据
		stmt, err := db.Prepare(`INSERT INTO ProblemInfo(problem,contactInfo,time) VALUES(?,?,?)`)
	aterr.CheckErr(err)


		res, err := stmt.Exec("HOUXIAOYUN", "TEST", "2017-06-24")
	aterr.CheckErr(err)
		id, err := res.LastInsertId()
	aterr.CheckErr(err)
		fmt.Println(id)
}
func updateData(db *sql.DB, id int) {


		//更新数据
		stmt, err := db.Prepare("update userinfo set username=? where uid=?")
	aterr.CheckErr(err)
		res, err := stmt.Exec("yanghaiUpdate", 23)
	aterr.CheckErr(err)
		affect, err := res.RowsAffected()
	aterr.CheckErr(err)
	fmt.Println(affect)
}

func queryData(db *sql.DB) {

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	data := []UserInfo{}
	aterr.CheckErr(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		aterr.CheckErr(err)
		ui := UserInfo{username,department,created}
		data = append(data,ui)
		fmt.Println(uid, username, department, created)
	}
	res,err :=json.Marshal(data)
	fmt.Println(string(res))
}

type UserInfo struct {
	UserName string
	DepartName string
	Created string
}
func CheckErr(err error) {
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}


type DbConfig struct{
	Name string
	Pwd string
	Ip string
	DbName string
}
