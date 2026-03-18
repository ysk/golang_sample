package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Hostname string
	URL      string
}

func (s Server) Check() string {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(s.URL)
	duration := time.Since(start).Round(time.Millisecond)

	if err != nil {
		errStr := err.Error()
		shortMsg := "Unknown Error"

		switch {
		case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline"):
			shortMsg = "Timeout"
		case strings.Contains(errStr, "connection refused"):
			shortMsg = "Refused"
		case strings.Contains(errStr, "forcibly closed") || strings.Contains(errStr, "wsarecv"):
			shortMsg = "ClosedByHost"
		case strings.Contains(errStr, "no such host"):
			shortMsg = "DNS Error"
		default:
			shortMsg = "Net Error"
		}

		// [NG] を赤色に
		return fmt.Sprintf("\x1b[31m[NG] %-20s | %-12s | %v\x1b[0m", s.Hostname, shortMsg, duration)
	}
	defer resp.Body.Close()

	return fmt.Sprintf("[OK] %-15s | Status: %d | Time: %v", s.Hostname, resp.StatusCode, duration)
}

// チャネルを使って並列チェックを実行する関数
func runCheck(servers []Server) {

	resultChan := make(chan string, len(servers))
	var wg sync.WaitGroup
	fmt.Println("--- 外部ファイル読み込み完了。チェック開始 ---")

	// --- 送信側 (Producer) ---
	for _, s := range servers {
		wg.Add(1)
		go func(target Server) {
			defer wg.Done()
			resultChan <- target.Check()
		}(s)
	}

	// --- 完了監視 ---
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// --- 受信側 (Consumer) ---
	successCount := 0
	for result := range resultChan {
		fmt.Println(result)
		if strings.Contains(result, "[OK]") {
			successCount++
		}
	}

	fmt.Println("---- 全台確認完了 ----")
	fmt.Printf("成功: %d / 合計: %d\n", successCount, len(servers))

}

func main() {

	file, err := os.Open("targets.txt")
	if err != nil {
		fmt.Printf("ファイルが開けませんでした: %v\n", err)
		return
	}
	defer file.Close()

	var servers []Server
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		servers = append(servers, Server{Hostname: parts[0], URL: parts[1]})
	}

	// 常時監視ループの設定
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// 最初の1回目を即座に実行
	fmt.Print("\033[H\033[2J")
	fmt.Printf("--- リアルタイム監視中: %s ---\n", time.Now().Format("15:04:05"))
	runCheck(servers)

	//常時監視ループ
	for {
		select {
		case <-ticker.C:
			fmt.Print("\033[H\033[2J")
			fmt.Printf("--- リアルタイム監視中: %s ---\n", time.Now().Format("15:04:05"))

			runCheck(servers)

			fmt.Println("==================================================")
		}
	}
}
