package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	uw = "username用户名:你想要设置的密码"
)

type authFileSrvHandler struct {
	http.Handler
}

func (f *authFileSrvHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	auth := r.Header.Get("Authorization")
	log.Println("got req:auth", auth)
	if auth == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="您必须输入用户名和密码"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("auth->", auth)
	split := strings.Split(auth, " ")
	if len(split) == 2 && split[0] == "Basic" {
		bytes, err := base64.StdEncoding.DecodeString(split[1])
		if err == nil && string(bytes) == uw {
			f.Handler.ServeHTTP(w, r)
			return
		}
	}
	w.Write([]byte("请联系相关人员获取用户名和密码！"))
}

/*
*

	通过这种方式修改文件服务器的根路径及端口

jelex@jelexxudeMacBook-Pro static-file-server % go run .\main.go -p 8888

	8888

exit status 2
jelex@jelexxudeMacBook-Pro static-file-server % go run .\main.go -r <absolute path>
jelex@jelexxudeMacBook-Pro static-file-server % go run main.go -r ~
*/
func main() {

	var rootPath, port string

	flag.StringVar(&rootPath, "r", "", "文件根目录")
	flag.StringVar(&port, "p", "6000", "文件服务器端口")

	flag.Parse()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(rootPath))

	mux.Handle("/", &authFileSrvHandler{fs})

	fmt.Println(rootPath, port)

	http.ListenAndServe(":"+port, mux)
}


