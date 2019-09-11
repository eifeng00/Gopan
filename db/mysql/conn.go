package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8")

	if err != nil {
		fmt.Printf("%s\n", err)
	}
	db.SetMaxOpenConns(1000)

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn : 返回数据库链接对象
func DBConn() *sql.DB {
	return db
}
