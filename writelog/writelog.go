package main

import (
	"fmt"
	"os"
)

// ファイルにメッセージを書き込む関数
func WriteLog(filename string, message string) error {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()
	f.WriteString(message + "\n")
	return nil
}

func main() {
	err := WriteLog("sample.log", "Hello Golang!")
	if err != nil {
		fmt.Println("ログ書き込み失敗:", err)
	} else {
		fmt.Println("ログ書き込み成功！")
	}
}
