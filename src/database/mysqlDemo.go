package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"io/ioutil"
)

func main() {
	config,err:= ReadDbConfig("/Users/yh/AndroidStudioProjects/GoWeb/config/dbConfig.json")
	if err != nil {
		fmt.Println("..............ttttt")
		return
	}

	db, err := sql.Open("mysql", config.Name+":"+config.Pwd+"@tcp("+config.Ip+")/"+config.DbName+"?charset=utf8")
	checkErr(err)

	insertData(db)

	//
	//queryData(db)
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
		stmt, err := db.Prepare("INSERT INTO userinfo(username,departname,created) VALUES(?,?,?)")
		checkErr(err)

		res, err := stmt.Exec("HOUXIAOYUN", "TEST", "2017-06-24")
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println(id)
}
func updateData(db *sql.DB, id int) {


		//更新数据
		stmt, err := db.Prepare("update userinfo set username=? where uid=?")
		checkErr(err)
		res, err := stmt.Exec("yanghaiUpdate", 23)
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
	fmt.Println(affect)
}

func queryData(db *sql.DB) {

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	data := []UserInfo{}
	checkErr(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
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
func checkErr(err error) {
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}

/**
	读取数据库的配置文件
 */
func ReadDbConfig(path string) (*dbConfig,error) {
	res,err := ioutil.ReadFile(path)
	if err!= nil {
		return nil,err
	}
	var config dbConfig

	err = json.Unmarshal(res,&config)
	if err != nil {
		return nil,err
	}
	return &config,nil
}
type dbConfig struct{
	Name string
	Pwd string
	Ip string
	DbName string
}
