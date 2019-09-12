package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// UserSingup : 通过用户名和密码完成user表的注册操作
func UserSingup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Printf("Failed to insert, err : %s", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Printf("Failed to insert , err : %s ", err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignin : 用户登陆检查
func UserSignin(username string, encpassword string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("username not fount：" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpassword {
		return true
	}
	return false
}

// UpdataToken : 刷新用户登陆的token
func UpdataToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}

	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	return true
}

// User : 与数据库中的User表对应
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LaseActiveAt string
	Status       int
}

//GetUserInfo : 用户信息查询
func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")

	if err != nil {
		fmt.Println(err.Error())
		return user, nil
	}

	defer stmt.Close()

	//执行查询操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
