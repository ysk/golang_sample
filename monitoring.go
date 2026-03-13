package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "https://google.com"
	logFile := "error_log.txt"

	fmt.Println("監視を開始しました。エラーは", logFile, "に記録されます。")

	for {
		resp, err := http.Get(url)
		now := time.Now().Format("2006-01-02 15:04:05")

		// HTTPステータスが200以外、または通信エラーの場合
		if err != nil || resp.StatusCode != http.StatusOK {
			msg := fmt.Sprintf("[%s] ERROR: %v (Status: %d)\n", now, err, resp.StatusCode)
			f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			f.WriteString(msg)
			f.Close()
			fmt.Print(msg)
		} else {
			fmt.Print(".")
			resp.Body.Close()
		}

		time.Sleep(10 * time.Second)
	}
}
