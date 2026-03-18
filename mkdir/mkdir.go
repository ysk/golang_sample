package main

import (
	"fmt"
	input "mkdir/utils"
	"os"
	"path/filepath"
	"time"
)

func main() {
	prefix, err := input.Input("作成するディレクトリの接頭後を入力してください")

	if err != nil || prefix == "" {
		fmt.Fprintln(os.Stderr, "有効な名前を入力してください")
		os.Exit(1)
	}

	now := time.Now()
	timestamp := now.Format("20060102_150405")

	dirName := prefix + "_" + timestamp

	err = os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "作成に失敗:", err)
		os.Exit(1)
	}

	fileName := "memo.txt"
	filePath := filepath.Join(dirName, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル作成失敗：", err)
		os.Exit(1)
	}

	defer file.Close()

	fmt.Printf("フォルダ '%s' とファイル '%s' を作成しました！\n", dirName, filePath)
}
