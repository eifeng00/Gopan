package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// docker run --name filestoresDB -p 3306:3306 -v /data/mysql/dates/:/var/lib/mysql -e MYSQL_DATABASE=fileserver -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7

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

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
