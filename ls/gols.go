package main

import (
	"fmt"
	"os"
)

func main() {
	//カレントディレクトリを開く
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Print("エラー", err)
		return
	}

	fmt.Println("---ファイル一覧---")
	for _, file := range files {
		prefix := " [FILE]"
		if file.IsDir() {
			prefix = " [DIR ]"
		}

		info, _ := file.Info()
		mode := info.Mode()
		size := info.Size()
		name := file.Name()
		mtime := info.ModTime().Format("2006/01/02 15:04:05")

		fmt.Printf("%s %s %s %8d %s\n", prefix, mtime, mode, size, name)
	}
}
