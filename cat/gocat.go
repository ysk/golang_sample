package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("使い方:gocat [ファイル名]")
	}
	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("エラー：ファイル '%s' が見つかりません\n", fileName)
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		fmt.Println("読み取り中にエラーが発生しました", err)
	}

}
