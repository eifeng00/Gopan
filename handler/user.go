package handler

import (
	dblayout "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwdsalt = "f**3#q2e3"
)

//SingupHandler : 处理用户注册请求
func SingupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	encpassword := util.Sha1([]byte(password + pwdsalt))
	suc := dblayout.UserSingup(username, encpassword)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}

}

//SinglnHandler : 登陆接口
func SinglnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		// http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	//1. 校验用户名及密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encpassword := util.Sha1([]byte(password + pwdsalt))
	pwdChecked := dblayout.UserSignin(username, encpassword)

	if !pwdChecked {
		w.Write([]byte("FAILD"))
		return
	}

	//2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayout.UpdataToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	//3. 登陆成功后重定向到首页
	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

// UserInfoHandler : 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	// 2. 验证token是否有效

	isvaildToken := IsTokenVaild(token)
	if !isvaildToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 3. 查询用户信息
	user, err := dblayout.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//GenToken : 生成40位的Token
func GenToken(username string) string {
	//40位字符md5(username + timestamp + token_slat) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//IsTokenVaild : token是否有效
func IsTokenVaild(token string) bool {
	// TODO: 判断token的时效性, 是否过期

	//TODO ： 从数据库表tbl_user_token查询username对应的token信息

	//TODO : 对比两个Token是否一致
	return true
}
