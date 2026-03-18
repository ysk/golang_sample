package main

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type CommandResult struct {
	Name   string
	Output string
	Error  error
}

// 実際にOSコマンドを叩く関数(ホントは非推奨)
func executeCommand(name string, arg ...string) CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, name, arg...).Output()

	return CommandResult{
		Name:   name,
		Output: string(out),
		Error:  err,
	}
}

func main() {
	commands := [][]string{
		{"vmstat", "1", "1"}, // CPU/メモリの概要
		{"free", "-m"},       // メモリ詳細 (MB単位)
		{"uptime"},           // ロードアベレージ
	}

	resultChan := make(chan CommandResult, len(commands))
	var wg sync.WaitGroup

	fmt.Println("--- 障害一次切り分け：リソース診断開始 ---")

	// 送信側 (Producer)
	for _, cmd := range commands {
		wg.Add(1)
		go func(c []string) {
			defer wg.Done()
			resultChan <- executeCommand(c[0], c[1:]...)
		}(cmd)
	}

	// --- 完了監視 ---
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 受信側 (Consumer)
	for res := range resultChan {
		fmt.Printf("\n=== [%s] の結果 ===\n", res.Name)
		if res.Error != nil {
			fmt.Printf("エラー発生: %v\n", res.Error)
			continue
		}
		fmt.Println(res.Output)
	}

	fmt.Println("\n--- 診断完了 ---")
}
