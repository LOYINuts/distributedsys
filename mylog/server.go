package mylog

import (
	"io"
	"log"
	"net/http"
	"os"
)

var slog *log.Logger

// 要为基本类型写方法，就新建一个类型
type filelog string

// 实现了io.Writer接口的Write方法
func (fl filelog) Write(data []byte) (int, error) {
	// 不存在则创建，只写入，往后附加
	fd, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	return fd.Write(data)
}

func Run(dest string) {
	// 参数:io.Writer,前缀，flag
	// 日期+时间
	slog = log.New(filelog(dest), "LOYINuts:", log.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	slog.Printf("%v\n", message)
}
