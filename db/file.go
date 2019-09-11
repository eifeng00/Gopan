package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

//OnFileUploadFinished : 文件上传完成，保存Meta信息到数据库
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	db := mydb.DBConn()
	stmt, err := db.Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Printf("Failded to prepare statement, err : %s", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rt, err := ret.RowsAffected(); nil == err {
		if rt <= 0 {
			fmt.Printf("File with hash : %s has been upload before", filehash)
			return false
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta : 从Mysql中获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	db := mydb.DBConn()
	stmt, err := db.Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1")

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		return nil, err
	}
	return &tfile, nil

}
