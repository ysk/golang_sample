package main

import (
	"fmt"
	hello "lesson/utils" // 「module名/フォルダ名」でインポート
	"os"
	"strconv"
)

func main() {
	// 1. 入力を受け取る (値 x と エラー err の2つを受け取る)
	x, err := hello.Input("type a price")
	if err != nil {
		fmt.Fprintln(os.Stderr, "入力エラー:", err)
		os.Exit(1)
	}

	// 2. 文字列を数値(int)に変換
	n, err := strconv.Atoi(x)
	if err != nil {
		fmt.Fprintln(os.Stderr, "数字を入力してください:", err)
		os.Exit(1)
	}

	// 3. 計算（1.1倍して整数に戻す）
	p := float64(n)
	fmt.Println("税込価格:", int(p*1.1))
}
