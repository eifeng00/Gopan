package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdataHandler)
	http.HandleFunc("/file/delete", handler.FileDelectHandler)

	http.HandleFunc("/user/signup", handler.SingupHandler)
	http.HandleFunc("/user/signin", handler.SinglnHandler)
	http.HandleFunc("/user/info", handler.UserInfoHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server , err : %s", err.Error())
	}
}
