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

	var lastRead, lastWrite uint64
	var lastRecv, lastSent uint64
	const interval = 3
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB = 1 << (10 * iota)
		GB = 1 << (10 * iota)
	)
	for {

		c, _ := cpu.Percent(time.Second, false)
		v, _ := mem.VirtualMemory()

		io, _ := disk.IOCounters()
		currentIO := io[""]

		nv, _ := net.IOCounters(true)
		var totalRecv, totalSent uint64
		for _, n := range nv {
			totalRecv += n.BytesRecv
			totalSent += n.BytesSent
		}

		var readSpeed, writeSpeed, recvSpeed, sentSpeed float64
		if lastRecv > 0 {
			readSpeed = float64(currentIO.ReadBytes-lastRead) / KB / interval
			writeSpeed = float64(currentIO.WriteBytes-lastWrite) / KB / interval
			recvSpeed = float64(totalRecv-lastRecv) / KB / interval
			sentSpeed = float64(totalSent-lastSent) / KB / interval
		}

		lastRead, lastWrite = currentIO.ReadBytes, currentIO.WriteBytes
		lastRecv, lastSent = totalRecv, totalSent

		fmt.Printf("[%s] CPU:%4.1f%% | Mem:%4.1f%% | IO-R:%6.1f | IO-W:%6.1f | Recv:%6.1f | Sent:%6.1f (KB/s)\n",
			time.Now().Format("15:04:05"),
			c[0],
			v.UsedPercent,
			readSpeed,
			writeSpeed,
			recvSpeed,
			sentSpeed,
		)

		time.Sleep(interval * time.Second)
	}
}
