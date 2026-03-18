package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func main() {
	fmt.Println("--- リアルタイム監視 (CPU/Mem/Disk/NW) ---")

	// 差分計算用の前回値保持変数
	var lastRead, lastWrite uint64
	var lastRecv, lastSent uint64

	// 定数の定義
	const interval = 3 // 計測の間隔（秒）
	const (
		_  = iota
		KB = 1 << (10 * iota) // 1024
		MB = 1 << (10 * iota) // 1048576
		GB = 1 << (10 * iota) // 1073741824
	)

	for {
		// CPU
		c, _ := cpu.Percent(time.Second, false)

		// Memory
		v, _ := mem.VirtualMemory()

		// Disk
		ioStats, _ := disk.IOCounters()
		var currentRead, currentWrite uint64
		for _, d := range ioStats {
			currentRead += d.ReadBytes
			currentWrite += d.WriteBytes
		}

		// Network
		nv, _ := net.IOCounters(true)
		var currentRecv, currentSent uint64
		for _, n := range nv {
			currentRecv += n.BytesRecv
			currentSent += n.BytesSent
		}

		var readSpeed, writeSpeed, recvSpeed, sentSpeed float64
		if lastRecv > 0 {
			readSpeed = float64(currentRead-lastRead) / float64(KB) / interval
			writeSpeed = float64(currentWrite-lastWrite) / float64(KB) / interval
			recvSpeed = float64(currentRecv-lastRecv) / float64(KB) / interval
			sentSpeed = float64(currentSent-lastSent) / float64(KB) / interval
		}

		lastRead, lastWrite = currentRead, currentWrite
		lastRecv, lastSent = currentRecv, currentSent

		// しきい値によって色を変える
		// CPUの色判定
		cpuColor := ""
		if c[0] > 80.0 {
			cpuColor = "\x1b[31m" //赤
		} else if c[0] > 50.0 {
			cpuColor = "\x1b[33m" //黄
		}

		// メモリの色判定
		memColor := ""
		if v.UsedPercent > 90.0 {
			memColor = "\x1b[31m" //赤
		} else if v.UsedPercent > 80.0 {
			memColor = "\x1b[33m" //黄
		}

		reset := "\x1b[0m" //デフォルト

		fmt.Printf("[%s] %sCPU:%4.1f%%%s | %sMem:%4.1f%%%s | IO-R:%7.1f | IO-W:%7.1f | Recv:%7.1f | Sent:%7.1f (KB/s)\n",
			time.Now().Format("15:04:05"),
			cpuColor, c[0], reset,
			memColor, v.UsedPercent, reset,
			readSpeed,
			writeSpeed,
			recvSpeed,
			sentSpeed,
		)

		// 指定した秒数待機
		time.Sleep((interval - 1) * time.Second)
	}
}
