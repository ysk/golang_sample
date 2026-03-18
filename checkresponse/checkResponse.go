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

type Checker interface {
	Check() string
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

	return fmt.Sprintf(
		"[OK] %-15s | Status: %d | Time: %v",
		s.Hostname, resp.StatusCode, duration)

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
			return
		}

		s := Server{
			Hostname: parts[0],
			URL:      parts[1],
		}

		servers = append(servers, s)

		// PHPの var_dump($servers) と同じやつ
		// spew.Dump(servers)
		// os.Exit(0)
	}

	var wg sync.WaitGroup
	fmt.Println("--- 外部ファイル読み込み完了。チェック開始 ---")

	for _, s := range servers {
		wg.Add(1)
		go func(target Server) {
			defer wg.Done()
			result := target.Check()
			fmt.Println(result)
		}(s)
	}

	wg.Wait()

	fmt.Println("---- 全台確認完了 ----")

}
