package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	fmt.Println("--- リアルタイムメモリ監視 ---")

	for {
		v, _ := mem.VirtualMemory()
		fmt.Printf("[%s] メモリ使用率: %.1f%%\n",
			time.Now().Format("15:04:05"),
			v.UsedPercent,
		)
		time.Sleep(2 * time.Second)
	}
}
